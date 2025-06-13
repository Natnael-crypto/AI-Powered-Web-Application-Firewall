import json
import os
from collections import defaultdict
from urllib.parse import unquote
import re

BAD_WORDS_BY_TYPE_DIR = "words/words_by_type/"
COMMON_BAD_WORDS_PATH = "words/bad_words.txt"

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


def load_badwords_by_type(badwords_dir):
    badwords_dict = {}
    for filename in os.listdir(badwords_dir):
        if filename.endswith(".txt"):
            category = filename.replace("badwords_", "").replace("badword_", "").replace(".txt", "")
            with open(os.path.join(badwords_dir, filename), 'r', encoding='utf-8') as f:
                badwords_dict[category] = [
                    line.strip().lower() for line in f.readlines() if line.strip() and not line.startswith('#')
                ]
    return badwords_dict


def load_common_badwords(file_path):
    with open(file_path, 'r', encoding='utf-8') as f:
        return [line.strip().lower() for line in f.readlines() if line.strip()]


# Load badwords from files
BADWORDS_BY_TYPE = load_badwords_by_type(BAD_WORDS_BY_TYPE_DIR)
COMMON_BADWORDS = load_common_badwords(COMMON_BAD_WORDS_PATH)


def count_injection_characters(text, feature_template):
    counts = feature_template.copy()
    for feature, chars in INJECTION_CHARACTERS.items():
        for char in chars:
            counts[feature] += text.count(char)
    return counts


def count_badwords_by_type(text):
    counts = defaultdict(int)
    decoded_text = unquote(text)
    words = re.split(r'\W+', decoded_text.lower())

    for category in BADWORDS_BY_TYPE.keys():
        counts[f"badword_{category}"] = 0

    for category, bad_words in BADWORDS_BY_TYPE.items():
        for word in bad_words:
            count = words.count(word.lower())
            if count > 0:
                counts[f"badword_{category}"] += count
    return counts

def count_common_badwords(text):
    decoded_text = unquote(text)
    words = re.split(r'\W+', decoded_text.lower())
    
    bad_word_count = 0
    for word in COMMON_BADWORDS:
        count = words.count(word.lower())
        if count > 0:
            bad_word_count += count
    return bad_word_count

def parse_request_for_type_prediction(request):
    """Extract features for type prediction model"""
    features = {key: 0 for key in INJECTION_CHARACTERS}
    for section in ["url", "headers", "body"]:
        text = request.get(section, "").replace("\n", " ").lower()
        features = count_injection_characters(text, features)

    full_text = f"{request.get('url', '')} {request.get('headers', '')} {request.get('body', '')}".lower()
    badword_counts = count_badwords_by_type(full_text)

    sorted_badword_counts = dict(sorted(badword_counts.items()))
    features.update(sorted_badword_counts)
    return features

def parse_request_for_anomaly_prediction(request):
    """Extract features for anomaly detection (flat badwords + char count + single badword field)"""

    features = {key: 0 for key in INJECTION_CHARACTERS}
    for section in ["url", "headers", "body"]:
        text = request.get(section, "")
        if isinstance(text, dict):
            text = json.dumps(text)
        text = str(text).replace("\n", " ").lower()
        features = count_injection_characters(text, features)

    full_text = f"{request.get('url', '')} {json.dumps(request.get('headers', ''))} {request.get('body', '')}".lower()
    features["badword"] = count_common_badwords(full_text)

    return features


