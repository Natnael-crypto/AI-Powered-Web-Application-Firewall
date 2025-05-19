import os
from enum import Enum

# Directory Configuration
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
MODELS_DIR = os.path.join(BASE_DIR, "ml_models")
DATABASE_DIR = os.path.join(BASE_DIR, "database")

# Database Configuration
DATABASE_FILE = "models.db"
DATABASE_PATH = os.path.join(DATABASE_DIR, DATABASE_FILE)

# Words Configuration
BAD_WORDS_DIR = os.path.join(BASE_DIR, "words")
BAD_WORDS_BY_TYPE_DIR = os.path.join(BAD_WORDS_DIR, "words_by_type")
COMMON_BAD_WORDS_FILE = "bad_words.txt"
COMMON_BAD_WORDS_PATH = os.path.join(BAD_WORDS_DIR, COMMON_BAD_WORDS_FILE)

# Backend API Endpoints
BACKEND_BASE_URL = "http://backend-service"
API_ENDPOINTS = {
    "changes": f"{BACKEND_BASE_URL}/ml/changes",
    "type_analysis": f"{BACKEND_BASE_URL}/ml/submit-analysis",
    "untrained_model": f"{BACKEND_BASE_URL}/ml/model/untrained",
    "model_results": f"{BACKEND_BASE_URL}/ml/model/results",
    "selected_model": f"{BACKEND_BASE_URL}/ml/model/selected",
}

# Model Configuration
# MODEL_FILE = "current_model.joblib"
# MODEL_PATH = os.path.join(MODELS_DIR, MODEL_FILE)
MODEL_FILE_EXTENSION = ".joblib"

# Model to ID Mappings
class ThreatType(Enum):
    COMMAND_INJECTION = 1
    DIRECTORY_TRAVERSAL = 2
    FILE_INCLUSION = 3
    LDAP_INJECTION = 4
    NOSQL_INJECTION = 5
    OPEN_REDIRECT = 6
    SQL_INJECTION = 7
    SSTI = 8
    XSS = 9
    XXE = 10

THREAT_TYPE_MAPPING = {
    ThreatType.COMMAND_INJECTION: "Command Injection",
    ThreatType.DIRECTORY_TRAVERSAL: "Directory Traversal",
    ThreatType.FILE_INCLUSION: "File Inclusion",
    ThreatType.LDAP_INJECTION: "LDAP Injection",
    ThreatType.NOSQL_INJECTION: "NoSQL Injection",
    ThreatType.OPEN_REDIRECT: "Open Redirect",
    ThreatType.SQL_INJECTION: "SQL Injection",
    ThreatType.SSTI: "Server-Side Template Injection",
    ThreatType.XSS: "Cross-Site Scripting",
    ThreatType.XXE: "XML External Entity"
}

# Ensure directories exist
os.makedirs(MODELS_DIR, exist_ok=True)
os.makedirs(DATABASE_DIR, exist_ok=True)



# Backend API Endpoints
BACKEND_BASE_URL = "http://backend-service"
API_ENDPOINTS = {
    "models": f"{BACKEND_BASE_URL}/ml/models",
    "selected_model": f"{BACKEND_BASE_URL}/ml/models/selected",
    "changes": f"{BACKEND_BASE_URL}/ml/changes",
    "type_analysis": f"{BACKEND_BASE_URL}/ml/analysis"
}

# Model Configuration

# Ensure directories exist
os.makedirs(MODELS_DIR, exist_ok=True)