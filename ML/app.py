from flask import Flask, request, jsonify
import threading
import joblib
import requests
import pandas as pd
import os
from dotenv import load_dotenv

load_dotenv()

from utils.parse_request import (
    parse_request_for_anomaly_prediction,
    parse_request_for_type_prediction
)

# General App Settings
DEBUG=os.getenv("DEBUG", "False").lower() == "true"

# API Endpoints
BACKEND_API_URL = os.getenv("BACKEND_API_URL")
BACKEND_API_ML_MODELS_PATH = os.getenv("BACKEND_API_ML_MODELS_PATH")
BACKEND_API_TYPE_ANALYSIS_PATH = os.getenv("BACKEND_API_TYPE_ANALYSIS_PATH")

# Model Paths
ANOMALY_PREDICTOR_MODEL_PATH = os.getenv("ANOMALY_PREDICTOR_MODEL_PATH")
TYPE_PREDICTOR_MODEL_PATH = os.getenv("TYPE_PREDICTOR_MODEL_PATH")

app = Flask(__name__)
app.config['DEBUG'] = DEBUG

# Global state
anomaly_detector_model = None
type_predictor_model = None

THREAT_TYPE_MAPPING = {
    1: "Command Injection",
    2: "Directory Traversal",
    3: "File Inclusion",
    4: "LDAP Injection",
    5: "NoSQL Injection",
    6: "Open Redirect",
    7: "SQL Injection",
    8: "Server-Side Template Injection",
    9: "Cross-Site Scripting (XSS)",
    10: "XML External Entity (XXE)",
}


'''
Anomaly Prediction
'''
# Load the anomaly predictor model
def load_anomaly_detector_model():
    global anomaly_detector_model
    try:
        anomaly_detector_model = joblib.load(ANOMALY_PREDICTOR_MODEL_PATH)
        app.logger.info("Anomaly predictor model loaded from {ANOMALY_PREDICTOR_MODEL_PATH}")
    except Exception as e:
        app.logger.error(f"Failed to load anomaly predictor model: {e}")


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
        app.logger.info("Type predictor model loaded successfully.")
    except Exception as e:
        app.logger.error(f"Failed to load type predictor model: {e}")


# Predict the type of the request using the type predictor model
def predict_type(data):
    """Predict the type of the request using the type predictor model"""
    try:
        features = parse_request_for_type_prediction(data)
        features_df = pd.DataFrame([features])
        
        prediction = type_predictor_model.predict(features_df)[0]
        proability = type_predictor_model.predict_proba(features_df)[0]
        
        return {
            "threat_type": THREAT_TYPE_MAPPING[prediction],
            "confidence": float(max(proability)),
        }
    except Exception as e:
        app.logger.error(f"Type prediction error: {e}")
        return None


# Predict the type and notify the backend API
def predict_and_notify(data, endpoint):
    """Predict the type and notify the backend API"""
    prediction = predict_type(data)
    if prediction is not None:
        try:
            request_id = data.get("request_id", "")
            response = requests.post(
                endpoint,
                json={
                    "request_id": request_id,
                    "threat_type": prediction["threat_type"],
                    "confidence": prediction["confidence"],
                },
                timeout=10,
            )
            if response.status_code == 200:
                app.logger.info(f"Type prediction sent successfully: {prediction}")
            else:
                app.logger.error(
                    f"Failed to send type prediction: {response.status_code}"
                )
        except Exception as e:
            app.logger.error(f"Error sending type prediction: {e}")

def init_system():
    load_anomaly_detector_model()
    load_type_predictor_model()


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


if __name__ == "__main__":
    init_system()
    app.run(host="0.0.0.0", port=8090)