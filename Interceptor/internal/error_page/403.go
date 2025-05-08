package error_page

import (
	"fmt"
	"log"
	"net/http"
)

func Send403Response(w http.ResponseWriter, RuleID int, RuleMessage string, Action string, Status int) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "text/html")

	page := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>403 Forbidden</title>
		<style>
			body {
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				background-color: #f9fafb;
				color: #1f2937;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
			}
			.container {
				text-align: center;
				background-color: #ffffff;
				padding: 40px;
				border-radius: 10px;
				box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
				max-width: 500px;
				width: 90%%;
			}
			h1 {
				font-size: 3rem;
				color: #dc2626;
				margin-bottom: 10px;
			}
			p {
				font-size: 1.1rem;
				margin: 15px 0;
			}
			.code-box {
				background-color: #f3f4f6;
				padding: 15px;
				border-radius: 8px;
				text-align: left;
				margin-top: 20px;
				font-family: monospace;
				font-size: 0.95rem;
			}
			a {
				color: #2563eb;
				text-decoration: none;
				font-weight: 500;
			}
			a:hover {
				text-decoration: underline;
			}
			@media (max-width: 600px) {
				h1 {
					font-size: 2.2rem;
				}
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>403 Forbidden</h1>
			<p>Access to this resource has been denied for security reasons.</p>
			<p>If you believe this is an error, please contact the administrator.</p>
			<div class="code-box">
				<strong>Rule ID:</strong> %d<br>
				<strong>Message:</strong> %s<br>
				<strong>Action Taken:</strong> %s<br>
				<strong>Status Code:</strong> %d
			</div>
			<p><a href="/">← Return to Home</a></p>
		</div>
	</body>
	</html>
	`, RuleID, RuleMessage, Action, Status)

	_, err := w.Write([]byte(page))
	if err != nil {
		log.Printf("Failed to write custom 403 page: %v", err)
	}
}
