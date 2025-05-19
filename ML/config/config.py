import os
from enum import Enum

# Directory Configuration
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
MODELS_DIR = os.path.join(BASE_DIR, "ml_models")
TYPE_PREDICTOR_DIR = os.path.join(MODELS_DIR, "type_predictor")
DATABASE_DIR = os.path.join(BASE_DIR, "database")
BAD_WORDS_DIR = os.path.join(BASE_DIR, "words")
BAD_WORDS_BY_TYPE_DIR = os.path.join(BAD_WORDS_DIR, "words_by_type")

# Database Configuration
DATABASE_FILE = "models.db"
DATABASE_PATH = os.path.join(DATABASE_DIR, DATABASE_FILE)

# Words Configuration
COMMON_BAD_WORDS_FILE = "bad_words.txt"
COMMON_BAD_WORDS_PATH = os.path.join(BAD_WORDS_DIR, COMMON_BAD_WORDS_FILE)

# Backend API Endpoints
BACKEND_BASE_URL = "http://backend-service"
API_ENDPOINTS = {
    "type_analysis": f"{BACKEND_BASE_URL}/ml/submit-analysis",
    "models": f"{BACKEND_BASE_URL}/ml/models",
}

# Model Configuration
MODEL_FILE_EXTENSION = ".joblib"
TYPE_PREDICTOR_FILE = "type_predictor.joblib"
TYPE_PREDICTOR_PATH = os.path.join(TYPE_PREDICTOR_DIR, TYPE_PREDICTOR_FILE)

# Threat Type Mapping
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

# Ensure directories exist
os.makedirs(MODELS_DIR, exist_ok=True)
os.makedirs(DATABASE_DIR, exist_ok=True)
os.makedirs(BAD_WORDS_DIR, exist_ok=True)
os.makedirs(BAD_WORDS_BY_TYPE_DIR, exist_ok=True)
os.makedirs(TYPE_PREDICTOR_DIR, exist_ok=True)