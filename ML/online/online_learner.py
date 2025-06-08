import os
import pandas as pd
import requests
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report, accuracy_score, f1_score, precision_score, recall_score
from sklearn.utils import resample
import joblib
from online.decide_and_update_model import decide_and_update_model
from online.get_next_random_forest_model_version import get_next_random_forest_model_version
from online.generate_new_log_csv import generate_csv
from dotenv import load_dotenv

load_dotenv(dotenv_path='../.env')
BACKEND_API_URL = os.getenv("BACKEND_API_URL")

def load_dataset(file_path):
    """Loads dataset from CSV file."""
    try:
        return pd.read_csv(file_path)
    except Exception as e:
        raise RuntimeError(f"Failed to load dataset from {file_path}: {e}")

def balance_new_data(new_df, abnormal_ratio=2):
    """
    Balances new data by:
    - Keeping all abnormal samples.
    - Sampling a fixed ratio of normal samples.
    """
    normal_df = new_df[new_df["label"] == 0]
    abnormal_df = new_df[new_df["label"] == 1]
    
    sampled_normal_df = resample(
        normal_df,
        replace=False,
        n_samples=min(len(normal_df), len(abnormal_df) * abnormal_ratio),
        random_state=42
    )
    
    return pd.concat([sampled_normal_df, abnormal_df])


def prepare_training_data(original_file, new_file=None, use_new_data=True):
    """
    Combines original dataset with optional new logs after balancing.
    """
    original_df = load_dataset(original_file)

    if new_file and use_new_data:
        new_df = load_dataset(new_file)
        balanced_new_df = balance_new_data(new_df)
        combined_df = pd.concat([original_df, balanced_new_df])
    else:
        combined_df = original_df

    combined_df = combined_df.sample(frac=1, random_state=42).reset_index(drop=True)
    return combined_df


def train_random_forest_model(df, model_path):
    """
    Trains a Random Forest model, saves it, and returns evaluation metrics.
    
    Returns:
        model: Trained RandomForestClassifier
        metrics: dict with accuracy, precision, recall, and f1 score
    """
    try:
        X = df.drop(columns=["label"])
        y = df["label"]

        X_train, X_test, y_train, y_test = train_test_split(
            X, y, test_size=0.2, random_state=42
        )

        rf_model = RandomForestClassifier(
            n_estimators=100,
            max_depth=10,
            min_samples_split=5,
            min_samples_leaf=2,
            max_features='sqrt',
            class_weight='balanced',
            bootstrap=True,
            oob_score=True,
            n_jobs=-1,
            random_state=42
        )

        rf_model.fit(X_train, y_train)

        y_pred = rf_model.predict(X_test)

        accuracy = accuracy_score(y_test, y_pred)
        precision = precision_score(y_test, y_pred, average='weighted', zero_division=0)
        recall = recall_score(y_test, y_pred, average='weighted', zero_division=0)
        f1 = f1_score(y_test, y_pred, average='weighted', zero_division=0)

        print(f"Accuracy: {accuracy:.4f}")
        print(f"Precision: {precision:.4f}")
        print(f"Recall: {recall:.4f}")
        print(f"F1 Score: {f1:.4f}")
        print("\nClassification Report:\n", classification_report(y_test, y_pred))

        joblib.dump(rf_model, model_path)
        print(f"Model saved to {model_path}")

        return {"accuracy":accuracy,"precision":precision,"recall":recall,"f1":f1}

    except Exception as e:
        raise RuntimeError(f"Training failed: {e}")

def get_model():
    path=f"{BACKEND_API_URL}ml/model/selected"
    response=requests.get(path,headers={"X-Service":"M"},verify=False)
    response.raise_for_status()
    res=response.json()
    return res['model']


# === MAIN EXECUTION ===
def onlineLearn():
    try:
        print("Getting training data...")
        generate_csv()
        original_data_file = "./data/original_data.csv"
        new_logs_file = "./data/new_log.csv"
        model_name = get_next_random_forest_model_version()
        model_save_path = f"../ml_models/{model_name}"

        print("Preparing training data...")
        training_df = prepare_training_data(original_data_file, new_logs_file, use_new_data=True)

        print("Training model...")
        metrics = train_random_forest_model(training_df, model_save_path)

        # Fetch current model info
        model = get_model()

        # Decide and possibly update
        result = decide_and_update_model(metrics, model, model_name)

        print("=====-=========")
        print(result)
        print("=====-=========")


        # Handle rejection
        if result.get("status") == "rejected":
            print("New model rejected. Removing saved file...")
            if os.path.exists(model_save_path):
                os.remove(model_save_path)
                print(f"Deleted: {model_save_path}")
            else:
                print(f"Model file not found for deletion: {model_save_path}")
        else:
            print("Model accepted and updated successfully.")

        return result

    except Exception as e:
        print(f"Error in online learning process: {e}")
        return {"status": "error", "message": str(e)}
    