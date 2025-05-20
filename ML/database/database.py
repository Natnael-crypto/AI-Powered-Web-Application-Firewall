import sqlite3
from sqlite3 import Error
from ML.config.config import API_ENDPOINTS, MODEL_FILE_EXTENSION, MODELS_DIR, DATABASE_PATH
import requests
import joblib
import os
import uuid
from datetime import datetime, timedelta

def create_connection():
    """Create a database connection to SQLite database"""
    conn = None
    try:
        conn = sqlite3.connect(DATABASE_PATH)
        conn.execute("PRAGMA foreign_keys = ON")
        return conn
    except Error as e:
        print(f"Error connecting to database: {e}")
    return conn

def initialize_database():
    """Initialize the database with required tables"""
    conn = create_connection()
    if conn is not None:
        try:
            cursor = conn.cursor()
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS models (
                    id TEXT PRIMARY KEY,
                    predecessor_id TEXT REFERENCES models(id) ON DELETE SET NULL,
                    number_requests_used INTEGER NOT NULL DEFAULT 0,
                    models_name TEXT NOT NULL UNIQUE,
                    accuracy REAL,
                    precision REAL,
                    recall REAL,
                    f1 REAL,
                    expected_accuracy REAL,
                    expected_precision REAL,
                    expected_recall REAL,
                    expected_f1 REAL,
                    selected BOOLEAN NOT NULL DEFAULT 0,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    last_trained_at TIMESTAMP,
                    train_every REAL
                )
            """)
            conn.commit()
        except Error as e:
            print(f"Error initializing database: {e}")
        finally:
            conn.close()

def fetch_and_store_models():
    """Fetch models from backend and store in database"""
    try:
        response = requests.get(API_ENDPOINTS["models"])
        if response.status_code == 200:
            models = response.json()
            conn = create_connection()
            if conn is not None:
                cursor = conn.cursor()
                for model in models:
                    cursor.execute("""
                        INSERT OR REPLACE INTO models (
                            id, models_name, accuracy, precision, recall, f1,
                            expected_accuracy, expected_precision, expected_recall, expected_f1,
                            selected, train_every, last_trained_at
                        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """, (
                        model["id"],
                        model["models_name"],
                        model.get("accuracy"),
                        model.get("precision"),
                        model.get("recall"),
                        model.get("f1"),
                        model.get("expected_accuracy"),
                        model.get("expected_precision"),
                        model.get("expected_recall"),
                        model.get("expected_f1"),
                        model.get("selected", False),
                        model.get("train_every"),
                        model.get("last_trained_at")
                    ))
                conn.commit()
                conn.close()
                return True
    except Exception as e:
        print(f"Error fetching models: {e}")
    return False

def get_selected_model():
    """Get the currently selected model"""
    conn = create_connection()
    if conn:
        try:
            cursor = conn.cursor()
            cursor.execute("""
                SELECT * FROM models WHERE selected = 1 ORDER BY updated_at DESC LIMIT 1
            """)
            return cursor.fetchone()
        except Error as e:
            print(f"Error getting selected model: {e}")
        finally:
            conn.close()
    return None

def update_model_file(model_id, model_name, model_data):
    """Save model file to disk"""
    try:
        model_path = os.path.join(MODELS_DIR, f"{model_name}{MODEL_FILE_EXTENSION}")
        with open(model_path, "wb") as f:
            f.write(model_data.content)
        return model_path
    except Exception as e:
        print(f"Error saving model file: {e}")
        return None

def should_train_model(selected):
    """Determine if the selected model needs retraining"""
    try:
        last_trained = selected[-2]
        interval = selected[-1]
        if interval is None:
            return False
        if last_trained is None:
            return True
        last_trained_time = datetime.fromisoformat(last_trained)
        return datetime.utcnow() >= last_trained_time + timedelta(minutes=interval)
    except Exception as e:
        print(f"Error checking training schedule: {e}")
        return False

def train_and_save_model():
    """Mock training logic. You should replace with real training."""
    try:
        new_model_id = str(uuid.uuid4())
        model_name = new_model_id  # model_name == file_name
        model_path = os.path.join(MODELS_DIR, f"{model_name}{MODEL_FILE_EXTENSION}")

        # Replace the next line with your actual model training
        dummy_model = {"predict": lambda x: [1]}  # Placeholder for a scikit-learn model
        joblib.dump(dummy_model, model_path)

        conn = create_connection()
        cursor = conn.cursor()
        cursor.execute("""
            INSERT INTO models (
                id, models_name, accuracy, precision, recall, f1,
                expected_accuracy, expected_precision, expected_recall, expected_f1,
                selected, train_every, last_trained_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            new_model_id, model_name,
            0.9, 0.9, 0.9, 0.9,  # dummy scores
            0.95, 0.95, 0.95, 0.95,
            False,
            10,  # train every 10 minutes (example)
            datetime.utcnow().isoformat()
        ))
        conn.commit()
        conn.close()

        return {"id": new_model_id, "model_name": model_name}
    except Exception as e:
        print(f"Error training and saving model: {e}")
        return None

def mark_predecessor(new_id, old_id):
    """Mark the current selected model as the predecessor of the new one"""
    conn = create_connection()
    if conn:
        try:
            cursor = conn.cursor()
            cursor.execute("""
                UPDATE models SET predecessor_id = ? WHERE id = ?
            """, (old_id, new_id))
            conn.commit()
        except Error as e:
            print(f"Error setting predecessor: {e}")
        finally:
            conn.close()

def update_model_version(new_model):
    """Mark new model as selected, unselect others"""
    conn = create_connection()
    if conn:
        try:
            cursor = conn.cursor()
            cursor.execute("UPDATE models SET selected = 0 WHERE selected = 1")
            cursor.execute("UPDATE models SET selected = 1 WHERE id = ?", (new_model["id"],))
            conn.commit()
        except Error as e:
            print(f"Error updating model version: {e}")
        finally:
            conn.close()

def delete_model_file(model_id):
    """Delete a model file and remove from DB"""
    conn = create_connection()
    if conn:
        try:
            cursor = conn.cursor()
            cursor.execute("SELECT models_name FROM models WHERE id = ?", (model_id,))
            row = cursor.fetchone()
            if row:
                file_name = f"{row[0]}{MODEL_FILE_EXTENSION}"
                path = os.path.join(MODELS_DIR, file_name)
                if os.path.exists(path):
                    os.remove(path)
            cursor.execute("DELETE FROM models WHERE id = ?", (model_id,))
            conn.commit()
        except Error as e:
            print(f"Error deleting model: {e}")
        finally:
            conn.close()

def discard_model(model):
    """Remove model from DB and disk if training was interrupted or replaced"""
    if model:
        delete_model_file(model["id"])
