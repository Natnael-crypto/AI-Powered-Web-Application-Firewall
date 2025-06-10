import os
import requests
from dotenv import load_dotenv

load_dotenv(dotenv_path='./.env')
BACKEND_API_URL = os.getenv("BACKEND_API_URL")

def decide_and_update_model(new_metrics, current_model_info,model_new_name):
    """
    Decide whether to accept or reject the new model based on both:
    - improvements over the current model
    - meeting expected minimum thresholds
    
    Args:
        new_metrics (dict): Metrics from the newly trained model.
        current_model_info (dict): Current model's metadata from the DB.

    Returns:
        dict: Decision result with message and optionally the response from the update.
    """
    try:
        # Extract old metrics
        old_accuracy = current_model_info.get('accuracy', 0)-1
        old_precision = current_model_info.get('precision', 0)-1
        old_recall = current_model_info.get('recall', 0)-1
        old_f1 = current_model_info.get('f1', 0)-1

        # Extract expected thresholds
        expected_accuracy = current_model_info.get('expected_accuracy', 0)
        expected_precision = current_model_info.get('expected_precision', 0)
        expected_recall = current_model_info.get('expected_recall', 0)
        expected_f1 = current_model_info.get('expected_f1', 0)

        # Extract new metrics and convert to %
        new_accuracy = new_metrics['accuracy'] * 100
        new_precision = new_metrics['precision'] * 100
        new_recall = new_metrics['recall'] * 100
        new_f1 = new_metrics['f1'] * 100

        # Conditions
        meets_expectation = (
            new_accuracy >= expected_accuracy and
            new_precision >= expected_precision and
            new_recall >= expected_recall and
            new_f1 >= expected_f1
        )

        improves_on_current = (
            new_accuracy >= old_accuracy and
            new_precision >= old_precision and
            new_recall >= old_recall and
            new_f1 >= old_f1
        )


        if meets_expectation and improves_on_current:
            # Accept model and update stats
            payload = {
                "id": current_model_info["id"],
                "accuracy": round(new_accuracy, 2),
                "precision": round(new_precision, 2),
                "recall": round(new_recall, 2),
                "f1": round(new_f1, 2),
                "model_type": model_new_name
            }

            path = f"{BACKEND_API_URL}ml/model/results"
            response = requests.post(path, json=payload, verify=False,headers={"X-Service":"M"})
            print("Saved the model")

            if response.status_code == 200:
                return {
                    "status": "accepted",
                    "message": "New model accepted and stats updated.",
                    "response": response.json()
                }
            else:
                return {
                    "status": "accepted",
                    "message": "New model accepted but failed to update stats.",
                    "response": response.text
                }
        else:
            reason = "Model did not meet expected thresholds." if not meets_expectation else "Model did not improve over current stats."
            return {
                "status": "rejected",
                "message": reason,
                "metrics_comparison": {
                    "accuracy": {"new": new_accuracy, "old": old_accuracy, "expected": expected_accuracy},
                    "precision": {"new": new_precision, "old": old_precision, "expected": expected_precision},
                    "recall": {"new": new_recall, "old": old_recall, "expected": expected_recall},
                    "f1": {"new": new_f1, "old": old_f1, "expected": expected_f1},
                }
            }

    except Exception as e:
        return {
            "status": "error",
            "message": str(e)
        }