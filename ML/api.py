from flask import Flask, request, jsonify
import re
import os
import numpy as np
import tensorflow as tf
from tensorflow.keras.models import load_model
from tensorflow.keras.preprocessing.sequence import pad_sequences

app = Flask(__name__)

# Load the pre-trained CNN model
curr_dir = os.path.dirname(os.path.realpath(__file__))
MODEL_PATH = os.path.join(curr_dir, "csic_model.keras")
model = load_model(MODEL_PATH)

def parse_headers(header_string):
    """
    Parses a raw header string into a dictionary.
    Example input: 'Accept:[text/html,application/xhtml+xml] User-Agent:[Mozilla/5.0]'

    Output: {'Accept': 'text/html,application/xhtml+xml', 'User-Agent': 'Mozilla/5.0'}
    """
    header_dict = {}
    matches = re.findall(r"(\S+):\[(.*?)\]", header_string)  # Match "Key:[Value]"

    for key, value in matches:
        header_dict[key] = value.replace(" ", "")  # Remove unnecessary spaces

    return header_dict


def preprocess_request(data):
    """
    Preprocesses a raw HTTP request dictionary into a fixed-length ASCII sequence for the model.
    
    request_data: Dictionary containing 'Url', 'Method', 'Headers' (as a string), and 'Body'
    """

    url = data.get("Url", "")
    method = data.get("Method", "")
    headers = data.get("Headers", "")
    body = data.get("Body", "")

    header_dict = parse_headers(headers)

    # order of fields: method, user_agent, pragma, cache_control, accept, Accept-Encoding,Accept-Charset, language, host, cookie, content-type, connection, content-length, body, url, http_version
    fields = [
        method,
        header_dict.get("User-Agent", ""),
        header_dict.get("Pragma", ""),
        header_dict.get("Cache-Control", ""),
        header_dict.get("Accept", ""),
        header_dict.get("Accept-Encoding", ""),
        header_dict.get("Accept-Charset", ""),
        header_dict.get("Accept-Language", "en"),
        header_dict.get("Host", ""),
        header_dict.get("Cookie", ""),
        header_dict.get("Content-Type", ""),
        "Connection: " + header_dict.get("Connection", ""),
        "Content-Length: " + header_dict.get("Content-Length", ""),
        body,
        url,
        header_dict.get("HTTP-Version", "HTTP/1.1")
    ]
   
    structured_data = ' '.join(fields)
    return structured_data

def encode_request(request_string):
    """Convert request string into a numerical format suitable for CNN processing."""
    # Simple encoding (convert characters to ASCII values)
    encoded = [ord(c) for c in request_string if ord(c) < 128]

    # Pad/truncate to fixed length
    padded_sequence = pad_sequences([encoded], maxlen=2000, padding='post', truncating='post')
    
    return padded_sequence[0]

@app.route("/predict", methods=["POST"])
def predict():
    try:
        json_data = request.get_json()
        if not json_data:
            return jsonify({"error": "Invalid JSON input"}), 400
        
        # Preprocess request
        formatted_request = preprocess_request(json_data)
        
        # Encode request
        encoded_request = encode_request(formatted_request)
        encoded_request = encoded_request.reshape(1, -1)
        
        
        # Make prediction
        prediction = model.predict(encoded_request)
        predicted_class = np.argmax(prediction, axis=1)[0]
        result = "Normal" if predicted_class == 0 else "Anomalous"
        
        return jsonify({"prediction": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":

    # Run the app on all network interfaces (0.0.0.0)
    app.run(host='0.0.0.0', debug=False)
