import os
import requests
import csv
from urllib.parse import unquote
import re
from collections import defaultdict
from dotenv import load_dotenv

load_dotenv(dotenv_path='../.env')

# Injection characters map
INJECTION_CHARACTERS = {
    "single_quote": ["'"],
    "double_quote": ["\""],
    "backtick": ["`"],
    "less_than": ["<"],
    "greater_than": [">"],
    "left_parenthesis": ["("],
    "right_parenthesis": [")"],
    "left_bracket": ["["],
    "right_bracket": ["]"],
    "left_brace": ["{{"],
    "right_brace": ["}}"],
    "dash": ["-"],
    "double_dash": ["--"],
    "hash": ["#"],
    "pipe": ["|"],
    "ampersand": ["&"],
    "dollar": ["$"],
    "percent": ["%"],
    "asterisk": ["*"],
    "exclamation_mark": ["!"],
    "equals": ["="],
    "logical_or": ["||"],
    "logical_and": ["&&"],
    "addition_operator": ["+"],
    "multiplication_operator": ["*"],
    "sql_comment_multi_line_open": ["/*"],
    "sql_comment_multi_line_close": ["*/"],
    "file_path_root": ["/"],
    "backslash": ["\\"],
    "colon": [":"],
    "comma": [","],
    "period": ["."],
    "caret": ["^"],
    "tilde": ["~"],
    "at_sign": ["@"],
    "carriage_return": ["\r"],
    "newline": ["\n"],
    "null_byte": ["\x00"],
    "space": [" "],
    "hexadecimal_prefix": ["0x"],
    "percent_encoded": ["%"],
}

# Path for bad words list
badwords_file = "./words/bad_words.txt"
BACKEND_API_URL = os.getenv("BACKEND_API_URL")

def read_bad_words(path):
    with open(path, 'r') as f:
        return set(line.strip().lower() for line in f if line.strip())

def count_injection_characters(text, injection_chars, features):
    for feature, chars in injection_chars.items():
        for char in chars:
            count = text.count(char)
            if count > 0:
                features[feature] += count
    return features

def count_bad_words(text, bad_words):
    decoded_text = unquote(text)
    words = re.split(r'\W+', decoded_text.lower())
    return sum(words.count(word) for word in bad_words)

def extract_features(request, injection_chars, bad_words):
    features = defaultdict(int)

    for field in ["url", "headers", "body"]:
        count_injection_characters(request.get(field, ""), injection_chars, features)

    full_text = f"{request.get('url', '')} {request.get('headers', '')} {request.get('body', '')}"
    features["badword"] = count_bad_words(full_text, bad_words)
    features["label"] = request.get("label", 0)
    return features

def fetch_requests():
    base_url =BACKEND_API_URL
    if not base_url:
        raise EnvironmentError("BACKEND_API_URL is not set in environment variables.")
    
    endpoint = f"{base_url}ml/requests"
    response = requests.get(endpoint,verify=False,headers={"X-Service":"M"})
    response.raise_for_status()
    return response.json()

def generate_csv():

    requests_data = fetch_requests()


    bad_words = read_bad_words(badwords_file)

    output_file = "./data/new_log.csv"
    os.makedirs(os.path.dirname(output_file), exist_ok=True)

    fieldnames = list(INJECTION_CHARACTERS.keys()) + ["badword", "label"]


    with open(output_file, "w", newline="") as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()


        for req in requests_data['requests']:
            features = extract_features(req, INJECTION_CHARACTERS, bad_words)
            writer.writerow({field: features.get(field, 0) for field in fieldnames})



