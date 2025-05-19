import os
from collections import defaultdict
from config.config import BAD_WORDS_BY_TYPE_DIR, COMMON_BAD_WORDS_PATH

# Shared character mappings
INJECTION_CHARACTERS = {
    "single_quote": ["'"],
    "double_quote": ['"'],
    "backtick": ["`"],
    "less_than": ["<"],
    "greater_than": [">"],
    "left_parenthesis": ["("],
    "right_parenthesis": [")"],
    "left_bracket": ["["],
    "right_bracket": ["]"],
    "left_brace": ["{"],
    "right_brace": ["}"],
    "dash": ["-"],
    "double_dash": ["--"],
    "hash": ["#"],
    "semicolon": [";"],
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
    "subtraction_operator": ["-"],
    "multiplication_operator": ["*"],
    "division_operator": ["/"],
    "sql_comment_single_line": ["--"],
    "sql_comment_multi_line_open": ["/*"],
    "sql_comment_multi_line_close": ["*/"],
    "file_path_root": ["/"],
    "backslash": ["\\"],
    "colon": [":"],
    "comma": [","],
    "period": ["."],
    "caret": ["^"],
    "tilde": ["~"],
    "question_mark": ["?"],
    "at_sign": ["@"],
    "vertical_tab": ["\v"],
    "tab": ["\t"],
    "carriage_return": ["\r"],
    "newline": ["\n"],
    "null_byte": ["\x00"],
    "space": [" "],
    "line_feed": ["\n"]
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
FLAT_BADWORDS = load_common_badwords(COMMON_BAD_WORDS_PATH)


def count_injection_characters(text, feature_template):
    counts = feature_template.copy()
    for feature, chars in INJECTION_CHARACTERS.items():
        for char in chars:
            counts[feature] += text.count(char)
    return counts


def count_badwords_by_type(text):
    text = text.lower()
    counts = defaultdict(int)
    for category, words in BADWORDS_BY_TYPE.items():
        for word in words:
            if word in text:
                counts[f"badword_{category}"] += text.count(word)
    return counts


def count_flat_badwords(text):
    text = text.lower()
    count = 0
    for word in FLAT_BADWORDS:
        if word in text:
            count += text.count(word)
    return count


def parse_request_for_type_prediction(request):
    """Extract features for type prediction model"""
    features = {key: 0 for key in INJECTION_CHARACTERS}
    for section in ["url", "headers", "body"]:
        text = request.get(section, "").replace("\n", " ").lower()
        features = count_injection_characters(text, features)

    full_text = f"{request.get('url', '')} {request.get('headers', '')} {request.get('body', '')}".lower()
    badword_counts = count_badwords_by_type(full_text)

    features.update(badword_counts)
    return features


def parse_request_for_anomaly_detection(request):
    """Extract features for anomaly detection (flat badwords + char count + single badword field)"""
    features = {key: 0 for key in INJECTION_CHARACTERS}
    for section in ["url", "headers", "body"]:
        text = request.get(section, "").replace("\n", " ").lower()
        features = count_injection_characters(text, features)

    full_text = f"{request.get('url', '')} {request.get('headers', '')} {request.get('body', '')}".lower()
    features["badword"] = count_flat_badwords(full_text)

    return features
