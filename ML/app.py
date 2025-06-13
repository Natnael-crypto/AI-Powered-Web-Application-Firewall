import glob
import re
import time
from flask import Flask, request, jsonify
import threading
import joblib
import requests
import pandas as pd
import os
from dotenv import load_dotenv
import json
import urllib3
from utils.parse_type_predction import extract_features
from online.online_learner import onlineLearn
from utils.parse_request import (
    parse_request_for_anomaly_prediction,
    parse_request_for_type_prediction
)
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

load_dotenv()

# General App Settings
DEBUG=os.getenv("DEBUG", "False").lower() == "true"

# API Endpoints
BACKEND_API_URL = os.getenv("BACKEND_API_URL")
BACKEND_API_ML_MODELS_PATH = "ml/models"
BACKEND_API_TYPE_ANALYSIS_PATH = "/ml/submit-analysis"
MODEL_CHANGES_ENDPOINT = "/ml/changes"

# Model Paths
ANOMALY_PREDICTOR_MODEL_PATH = "ml_models/random_forest_model_v.0.2.3.U.pkl"
TYPE_PREDICTOR_MODEL_PATH = "ml_models/random_forest_model_type.pkl"

app = Flask(__name__)
app.config['DEBUG'] = DEBUG

# Global state
anomaly_detector_model = None
type_predictor_model = None

latest_model_name = None
train_every_interval = None

prediction_queue = []
queue_lock = threading.Lock()
MAX_QUEUE_SIZE = 1000


SELECTED_MODEL_ENDPOINT = "ml/model/selected"
model_change_watcher_started = False

THREAT_TYPE_MAPPING = {
    1: "Command Injection",
    2: "Path Traversal",
    3: "File Inclusion",
    4: "LDAP Injection",
    5: "NoSQL Injection",
    6: "SQL Injection",
    7: "Server-Side Template Injection (SSTI)",
    8: "Cross-Site Scripting (XSS)",
    9: "XML External Entity (XXE)",
}


'''
Anomaly Prediction
'''

attack_types = {
    "command": 1, "directory_traversal": 2, "file_inclusion": 3, "ldap": 4,
    "nosql": 5, "sql": 6, "sst": 7, "xss": 8, "xxe": 9
}

# Flush queue to backend
def flush_queue():
    with queue_lock:
        if not prediction_queue:
            return
        to_send = prediction_queue[:]
        prediction_queue.clear()
    try:
        response = requests.post(
            BACKEND_ENDPOINT,
            headers={"X-Service": "M"},
            json=to_send,
            timeout=10,
            verify=False
        )
        if response.status_code == 200:
            app.logger.info(f"Batch sent successfully: {len(to_send)} predictions")
        else:
            app.logger.error(f"Failed to send batch: {response.status_code}")
    except Exception as e:
        app.logger.error(f"Error sending predictions: {e}")


def watch_model_changes(poll_interval=120):
    """
    Poll the backend /ml/changes endpoint every `poll_interval` seconds to detect changes.
    If model_setting_updated is True, reinitialize the system.
    """
    def run():
        nonlocal_last_flag = {"model_setting_updated": False}

        while True:
            try:
                response = requests.get(f"{BACKEND_API_URL}{MODEL_CHANGES_ENDPOINT}", verify=False, headers={"X-Service": "M"},)
                response.raise_for_status()

                data = response.json()
                updated_flag = data.get("model_setting_updated", False)

                if updated_flag:
                    app.logger.info("Model settings updated. Reinitializing system...")
                    init_system()

            except Exception as e:
                app.logger.error(f"Error checking for model changes: {e}")

            time.sleep(poll_interval)

    thread = threading.Thread(target=run, daemon=True)
    thread.start()


def get_latest_anomaly_model_path():
    """
    Fetch the selected model name from the API and return the latest matching model file path.
    """
    try:
        # Fetch the selected model info
        response = requests.get(f"{BACKEND_API_URL}{SELECTED_MODEL_ENDPOINT}", verify=False,headers={"X-Service":"M"})
        response.raise_for_status()
        model = response.json()

        model_type = model['model']["model_type"]
        if not model_type:
            raise ValueError("Model type not found in API response")

        # Pattern for model files e.g., random_forest_model_v.0.2.3.U.pkl
        model_pattern = f"./ml_models/{model_type}"
        model_files = glob.glob(model_pattern)


        if not model_files:
            raise FileNotFoundError(f"No model file found for pattern: {model_pattern}")

        # Extract version number from file name
        def extract_version_num(fname):
            match = re.search(r"v\.0\.2\.(\d+)\.U\.pkl", fname)
            return int(match.group(1)) if match else -1

        latest_model = sorted(model_files, key=extract_version_num)[-1]
        return latest_model,model['model']

    except requests.RequestException as e:
        raise RuntimeError(f"Failed to fetch model info from backend: {e}")
    except Exception as e:
        raise RuntimeError(f"Failed to determine latest model path: {e}")


def load_latest_anomaly_model():
    """
    Load the most recent anomaly model and get train_every from it
    """
    global anomaly_detector_model, latest_model_name, train_every_interval

    model_path,ai_models = get_latest_anomaly_model_path()
    anomaly_detector_model = joblib.load(model_path)
    latest_model_name = os.path.basename(model_path)
    app.logger.info(f"Loaded anomaly model: {latest_model_name}")

    # Load train_every from model object (or fallback metadata strategy)
    try:
        train_every_interval =ai_models["train_every"]
        if not train_every_interval:
            raise ValueError("Model missing `train_every` attribute")

        app.logger.info(f"train_every set to {train_every_interval} ms")
    except Exception as e:
        app.logger.error(f"Failed to read train_every from model: {e}")
        train_every_interval = None


# Prediction utility function
def predict_anomality_with_model(input_data):
    global anomaly_detector_model
    try:
        # Convert input to DataFrame if necessary
        if isinstance(input_data, dict):
            input_data = pd.DataFrame([input_data])
        elif isinstance(input_data, list):
            input_data = pd.DataFrame(input_data)

        # Ensure model is loaded
        if anomaly_detector_model is None:
            raise Exception("Model not loaded")

        # Predict
        predictions = anomaly_detector_model.predict(input_data)
        proba = anomaly_detector_model.predict_proba(input_data)
        return predictions[0],proba[0]  # Assuming single prediction
    except Exception as e:
        app.logger.error(f"Anomaly prediction failed: {e}")
        return None

'''
Type Prediction
'''
# Load the type predictor model
def load_type_predictor_model():
    """Load the type predictor model"""
    global type_predictor_model
    try:
        type_predictor_model = joblib.load(TYPE_PREDICTOR_MODEL_PATH)
        app.logger.info(f"Type predictor model loaded successfully: {type(type_predictor_model)}")
    except Exception as e:
        app.logger.error(f"Failed to load type predictor model: {e}")


# Predict the type of the request using the type predictor model
def predict_type(data):
    """Predict the type of the request using the type predictor model"""
    try:
        features = extract_features(data)
        features_df = pd.DataFrame([features])
        
        prediction = type_predictor_model.predict(features_df)[0]
        probability = type_predictor_model.predict_proba(features_df)[0]
        return {
            "threat_type": THREAT_TYPE_MAPPING[prediction],
            "confidence": float(max(probability)),
        }
    except Exception as e:
        app.logger.error(f"Type prediction error: {e}")
        return None

def add_result_to_queue(result):
    global last_flush_time
    with queue_lock:
        prediction_queue.append(result)
        queue_size = len(prediction_queue)

    if queue_size >= MAX_QUEUE_SIZE:
        flush_queue()

def predict_and_notify(data, endpoint):
    global BACKEND_ENDPOINT
    BACKEND_ENDPOINT = endpoint
    prediction = predict_type(data)
    if prediction is not None:
        request_id = data.get("request_id", "")
        result = {
            "request_id": request_id,
            "threat_type": prediction["threat_type"]
        }
        add_result_to_queue(result)




def schedule_online_learning(interval_ms):
    """
    Run onlineLearn() every `interval_ms` in a background thread
    """
    def run():
        while True:
            time.sleep(interval_ms / 1000.0)
            try:
                app.logger.info("Running scheduled online learning...")
                onlineLearn()
            except Exception as e:
                app.logger.error(f"Online learning error: {e}")

    thread = threading.Thread(target=run, daemon=True)
    thread.start()



@app.route("/predict", methods=["POST"])
def analyze_request():
    try:
        request_json = request.json
        if not request_json:
            return jsonify({"success": False, "error": "No JSON received"}), 400
        features = parse_request_for_anomaly_prediction(request_json)

        if anomaly_detector_model:
            prediction,prob = predict_anomality_with_model(features)
            if prediction is None:
                return jsonify({"success": False, "error": "Prediction failed"}), 500

            # # Send type analysis
            if prediction == 1:
                threading.Thread(
                    target=predict_and_notify,
                    args=(request_json, f"{BACKEND_API_URL}{BACKEND_API_TYPE_ANALYSIS_PATH}"),
                    daemon=True,
                ).start()

            print(prob)
            result = "Normal" if prediction == 0 else "Anomaly"
            return jsonify({"success": True, "prediction": result,"Normal":prob[0],"Anomaly":prob[1]}), 200
            

        return jsonify({"error": "No model loaded"}), 503

    except Exception as e:
        app.logger.error(f"Prediction error: {e}")
        return jsonify({"error": str(e)}), 500

# Error handlers
@app.errorhandler(404)
def not_found_error(e):
    return jsonify({"success": False, "error": "Not found"}), 404

@app.errorhandler(500)
def internal_error(e):
    return jsonify({"success": False, "error": "Internal server error"}), 500

def init_system(start_watcher=False):
    global model_change_watcher_started
    load_latest_anomaly_model()
    load_type_predictor_model()

    if train_every_interval:
        schedule_online_learning(train_every_interval)
    else:
        app.logger.warning("train_every is not defined; online learning will not run.")

    if start_watcher and not model_change_watcher_started:
        watch_model_changes(120)
        model_change_watcher_started = True

if __name__ == "__main__":
    init_system(start_watcher=True)
    app.run(host="0.0.0.0", port=8090, ssl_context=("./certs/cert.pem", "./certs/key.pem"))