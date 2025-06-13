import os
import re
from urllib.parse import unquote

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
    "asterisk": ["*"],
    "exclamation_mark": ["!"],
    "equals": ["="],
    "logical_or": ["||"],
    "logical_and": ["&&"],
    "addition_operator": ["+"],
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

attack_types = {
    "command": 1, "directory_traversal": 2, "file_inclusion": 3, "ldap": 4,
    "nosql": 5, "sql": 6, "sst": 7, "xss": 8, "xxe": 9
}

def read_bad_words(file_path):
    with open(file_path, 'r', encoding='utf-8') as f:
        return set(line.strip().lower() for line in f if line.strip())

def read_all_badwords(badwords_dir):
    badwords_dict = {}
    for attack in attack_types:
        path = os.path.join(badwords_dir, f"{attack}.txt")
        if os.path.exists(path):
            badwords_dict[attack] = read_bad_words(path)
        else:
            print(f"[!] Missing: {attack}.txt")
            badwords_dict[attack] = set()
    return badwords_dict

BADWORDDICT=read_all_badwords('words/words_by_type')

def count_injection_characters(text, injection_chars, features):
    for feature, chars in injection_chars.items():
        for char in chars:
            features[feature] += text.count(char)
    return features

def count_bad_words(text, bad_words):
    decoded_text = unquote(text)
    words = re.split(r'\W+', decoded_text.lower())
    count=0
    for word in bad_words:
        if words.count(word.lower())>0:
            count+=words.count(word.lower())
    return count



def extract_features(request):
    features = {feature: 0 for feature in INJECTION_CHARACTERS}
    full_text = request["url"] + " " +request["headers"] + " " + request["body"]
    
    # Count injection characters
    for part in [request["url"], request["headers"], request["body"]]:
        count_injection_characters(part, INJECTION_CHARACTERS, features)

    # Count bad words per attack type
    for attack_name, wordlist in BADWORDDICT.items():
        features[f"{attack_name}_badword"] = count_bad_words(full_text, wordlist)

    return features



