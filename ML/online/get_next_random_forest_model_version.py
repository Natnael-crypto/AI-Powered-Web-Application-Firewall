import os
import re

def get_next_random_forest_model_version(directory="../ml_models/"):
    """
    Searches for the current in-use random forest model and returns a new filename with incremented version.

    Returns:
        str: New versioned filename (e.g., random_forest_model_v.0.2.2.U.pkl)
    """
    try:
        pattern = re.compile(r"random_forest_model_v\.(\d+)\.(\d+)\.(\d+)\.U\.pkl")
        latest_patch = -1
        current_file = ""

        # List files in directory
        for filename in os.listdir(directory):
            match = pattern.match(filename)
            if match:
                major, minor, patch = map(int, match.groups())
                if patch > latest_patch:
                    latest_patch = patch
                    current_file = filename

        if not current_file:
            raise FileNotFoundError("No in-use model file found in expected format.")

        # Extract and increment patch
        major, minor, patch = map(int, pattern.match(current_file).groups())
        next_patch = patch + 1

        # Build next version filename
        next_filename = f"random_forest_model_v.{major}.{minor}.{next_patch}.U.pkl"

        return next_filename

    except Exception as e:
        raise RuntimeError(f"Failed to determine next model version: {e}")
