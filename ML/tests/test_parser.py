import pytest
from utils.parse_request import parse_request_for_anomaly_detection, parse_request_for_type_prediction
from collections import defaultdict


# Sample test config and badwords
mock_config = {
    "badwords_file": "mock_badwords.txt",
    "feature_columns": [
        "single_quote", "double_quote", "semicolon", "less_than",
        "greater_than", "badword_sqli", "badword_xss"
    ],
    "injection_characters": {
        "single_quote": ["'"],
        "double_quote": ['"'],
        "semicolon": [";"],
        "less_than": ["<"],
        "greater_than": [">"]
    }
}

sample_request = {
    "url": "/search?query=select+*+from+users",
    "headers": "Host: example.com\nUser-Agent: curl/7.68.0",
    "body": "<script>alert(1)</script>"
}

def test_extract_type_features():
    # features = extract_type_features(sample_request, mock_config, mock_badwords)
    features = parse_request_for_type_prediction(sample_request)

    # Check all expected features are present
    assert isinstance(features, dict)
    for col in mock_config["feature_columns"]:
        assert col in features

    # Check character counts
    assert features["single_quote"] == 0
    assert features["less_than"] >= 1
    assert features["badword_sqli"] >= 1
    assert features["badword_xss"] >= 1

def test_extract_anomaly_features():
    # features = extract_anomaly_features(sample_request, mock_config, mock_badwords)
    features = parse_request_for_anomaly_detection(sample_request)

    assert isinstance(features, dict)
    # Must contain all characters from injection_characters
    for key in mock_config["injection_characters"]:
        assert key in features

    assert features["semicolon"] == 0
    assert features["less_than"] >= 1
    assert features["greater_than"] >= 1
    assert "badword" in features
    assert features["badword"] >= 2
