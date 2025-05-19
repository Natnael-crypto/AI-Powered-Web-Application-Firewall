from flask import Flask, request, jsonify
import threading
import joblib
import os
import time
import requests

from config import API_ENDPOINTS, MODELS_DIR, MODEL_FILE_EXTENSION
from database import (
    initialize_database,
    fetch_and_store_models,
    get_selected_model,
    update_model_file,
    should_train_model,
    train_and_save_model,
    update_model_version,
    delete_model_file,
    discard_model,
    mark_predecessor,
)
from type_predictor import TypePredictor
from utils import extract_features

app = Flask(__name__)

# Global state
current_model = None
current_model_id = None
training_thread = None
training_lock = threading.Lock()
stop_training_flag = threading.Event()
type_predictor = TypePredictor()

def load_selected_model():
    """Load selected model from disk into global state"""
    global current_model, current_model_id
    selected = get_selected_model()
    if not selected:
        current_model = None
        current_model_id = None
        return

    model_id, _, model_name = selected
    model_path = os.path.join(MODELS_DIR, f"{model_name}{MODEL_FILE_EXTENSION}")

    if os.path.exists(model_path):
        current_model = joblib.load(model_path)
        current_model.model_id = model_id
        current_model.model_name = model_name
        current_model_id = model_id
        app.logger.info(f"Loaded model {model_id} ({model_name})")
    else:
        app.logger.warning(f"Model file not found: {model_path}")
        current_model = None
        current_model_id = None

def periodic_model_watcher():
    """Watches for remote changes and triggers training if necessary"""
    while True:
        try:
            if fetch_and_store_models():
                prev_model_id = current_model_id
                load_selected_model()
                if current_model_id != prev_model_id:
                    app.logger.info("Switched to new selected model")

            ensure_model_file_present()

            selected = get_selected_model()
            if selected and should_train_model(selected):
                start_training_if_appropriate(selected)

        except Exception as e:
            app.logger.error(f"Error during model watcher loop: {e}")

        time.sleep(60)

def ensure_model_file_present():
    """Fetch model file if missing"""
    selected = get_selected_model()
    if not selected:
        return

    model_id, _, model_name = selected
    model_file_path = os.path.join(MODELS_DIR, f"{model_name}{MODEL_FILE_EXTENSION}")
    if not os.path.exists(model_file_path):
        resp = requests.get(f"{API_ENDPOINTS['models']}/{model_id}/file")
        if resp.status_code == 200:
            update_model_file(model_id, model_name, resp)

def start_training_if_appropriate(selected):
    """Start model training in a thread if not already running"""
    global training_thread

    if training_thread and training_thread.is_alive():
        app.logger.info("Training already in progress. Skipping new training.")
        return

    training_thread = threading.Thread(target=train_model_flow, args=(selected,), daemon=True)
    training_thread.start()

def train_model_flow(selected):
    """Full training flow with model versioning and rollback checks"""
    model_id, _, model_name = selected
    old_model_id = model_id

    with training_lock:
        stop_training_flag.clear()
        app.logger.info("Starting model training...")

        try:
            new_model = train_and_save_model()
            if stop_training_flag.is_set():
                discard_model(new_model)
                app.logger.info("Training interrupted, new model discarded.")
                return

            current = get_selected_model()
            if current and current[0] != old_model_id:
                discard_model(new_model)
                app.logger.info("User changed selection mid-training. Discarding model.")
                return

            mark_predecessor(new_model['id'], old_model_id)
            delete_model_file(old_model_id)
            update_model_version(new_model)
            load_selected_model()
            notify_api_of_update()

        except Exception as e:
            app.logger.error(f"Training failed: {e}")

def notify_api_of_update():
    """Notify remote API of model update"""
    try:
        requests.post(API_ENDPOINTS['models'], json={"update": True})
    except Exception as e:
        app.logger.error(f"Failed to notify API of model update: {e}")

@app.route('/analyze', methods=['POST'])
def analyze_request():
    try:
        data = request.json
        features = extract_features(data)

        if current_model:
            prediction = current_model.predict([features])[0]
            if prediction == 1:
                threading.Thread(
                    target=type_predictor.predict_and_notify,
                    args=(data, API_ENDPOINTS['type_analysis']),
                    daemon=True
                ).start()

            return jsonify({
                "malicious": bool(prediction),
                "model": current_model.model_name
            })

        return jsonify({"error": "No model loaded"}), 503

    except Exception as e:
        app.logger.error(f"Prediction error: {e}")
        return jsonify({"error": str(e)}), 500

@app.before_first_request
def init_system():
    initialize_database()
    fetch_and_store_models()
    load_selected_model()
    threading.Thread(target=periodic_model_watcher, daemon=True).start()

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
