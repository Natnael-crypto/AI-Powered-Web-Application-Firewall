package error_page

import (
	"log"
	"net/http"
)

// Send403Response sends a styled 403 Forbidden page with a custom message.
func Send403Response(w http.ResponseWriter, RuleID int, RuleMessage string, Action string, Status int) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>403 Forbidden</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					text-align: center;
					background-color: #f4f4f4;
					color: #333;
					margin: 0;
					padding: 0;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
				}
				.container {
					max-width: 600px;
					background: #fff;
					padding: 30px;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
					border-radius: 8px;
				}
				h1 {
					font-size: 3em;
					color: #e74c3c;
				}
				p {
					font-size: 1.2em;
					margin: 20px 0;
				}
				a {
					color: #3498db;
					text-decoration: none;
				}
				a:hover {
					text-decoration: underline;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>403 Forbidden</h1>
				<p>Your request has been blocked.</p>
				<p><a href="/">Go back to Home</a></p>
			</div>
		</body>
		</html>
	`))
	if err != nil {
		log.Printf("Failed to write custom 403 page: %v", err)
	}
}
