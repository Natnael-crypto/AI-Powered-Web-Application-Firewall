from flask import Flask, render_template_string, request
from flask_socketio import SocketIO, send

app = Flask(__name__)
socketio = SocketIO(app)

html = """
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Vulnerable WebSocket & HTTP Form Test</title>
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
        <style>
            body {
                background-color: #f8f9fa;
                font-family: 'Arial', sans-serif;
            }
            h1 {
                color: #007bff;
                font-weight: 600;
            }
            .container {
                margin-top: 50px;
            }
            .form-control {
                margin-bottom: 10px;
            }
            .alert-info {
                font-weight: bold;
            }
            .list-group-item {
                padding: 10px;
            }
            .message-box {
                margin-top: 20px;
                padding: 10px;
                background-color: #e9ecef;
                border-radius: 5px;
                border: 1px solid #ced4da;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1 class="text-center mb-4">Vulnerable WebSocket and HTTP Form Test</h1>

            <!-- Normal HTTP Form -->
            <div class="row">
                <div class="col-md-6 mx-auto message-box">
                    <h3 class="text-center mb-4">Submit via HTTP</h3>
                    <form method="POST" action="/submit">
                        <div class="form-group">
                            <label for="messageText">Enter a message:</label>
                            <input type="text" id="messageText" name="message" class="form-control" placeholder="Type a message..." required>
                        </div>
                        <button type="submit" class="btn btn-primary btn-block">Submit</button>
                    </form>
                    <!-- XSS Vulnerability: directly inserting user input into HTML -->
                    {% if message %}
                        <div class="alert alert-info mt-3">
                            <strong>Submitted message:</strong> {{ message|safe }}
                        </div>
                    {% endif %}
                </div>
            </div>

            <!-- WebSocket Section -->
            <div class="row mt-5">
                <div class="col-md-6 mx-auto message-box">
                    <h3 class="text-center mb-4">WebSocket Test</h3>
                    <input type="text" id="wsMessageText" class="form-control" placeholder="Type a WebSocket message..." required>
                    <button onclick="sendMessage()" class="btn btn-success btn-block mt-3">Send via WebSocket</button>
                    <ul id="messages" class="list-group mt-4">
                    </ul>
                </div>
            </div>
        </div>

        <script>
            let ws = new WebSocket("ws://waf.local/socket");

            ws.onmessage = function(event) {
                let messages = document.getElementById('messages');
                let message = document.createElement('li');
                message.textContent = 'Server: ' + event.data;
                message.classList.add('list-group-item', 'bg-light');
                messages.appendChild(message);
            };

            function sendMessage() {
                let input = document.getElementById("wsMessageText");
                let message = input.value;
                if (message.trim() !== '') {
                    ws.send(message);
                    let messages = document.getElementById('messages');
                    let messageItem = document.createElement('li');
                    messageItem.textContent = 'You: ' + message;
                    messageItem.classList.add('list-group-item', 'bg-light');
                    messages.appendChild(messageItem);
                    input.value = '';
                }
            }
        </script>
        <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    </body>
</html>
"""

@app.route('/')
def index():
    return render_template_string(html, message=None)

@app.route('/submit', methods=['POST'])
def submit():
    # Vulnerable to XSS and SSTI: we use `request.form.get` without sanitization
    message = request.form.get('message')
    # No filtering or escaping, allowing potential XSS
    return render_template_string(html, message=message)

@socketio.on('message')
def handle_message(message):
    print(f"Received message: {message}")
    send(f"Echo: {message}")

if __name__ == '__main__':
    socketio.run(app, host='0.0.0.0', port=5000)
