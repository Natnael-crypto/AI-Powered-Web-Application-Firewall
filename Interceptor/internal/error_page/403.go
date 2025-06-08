package error_page

import (
	"fmt"
	"log"
	"net/http"
)

func Send403Response(w http.ResponseWriter, RequestID string) {
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
				margin: 0;
				padding: 0;
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
				background-color: #f4f6f8;
				color: #333;
				display: flex;
				flex-direction: column;
				justify-content: center;
				align-items: center;
				height: 100vh;
			}
			.wrapper {
				text-align: center;
				max-width: 600px;
				padding: 30px;
			}
			h1 {
				font-size: 72px;
				margin-bottom: 10px;
				color: #dc2626;
			}
			h2 {
				font-size: 24px;
				margin-bottom: 20px;
			}
			p {
				font-size: 16px;
				margin: 10px 0;
			}
			.diagnostic {
				background-color: #f0f0f0;
				border-left: 4px solid #dc2626;
				padding: 15px;
				text-align: left;
				margin-top: 25px;
				font-family: monospace;
				font-size: 14px;
				border-radius: 5px;
			}
			.footer {
				margin-top: 40px;
				font-size: 13px;
				color: #888;
			}
			a {
				color: #2563eb;
				text-decoration: none;
			}
			a:hover {
				text-decoration: underline;
			}
		</style>
	</head>
	<body>
		<div class="wrapper">
			<h1>403</h1>
			<h2>Access Denied</h2>
			<p>Your request has been blocked for security reasons.</p>
			<p>If you believe this was done in error, please contact the site administrator with the information below.</p>
			<div class="diagnostic">
				<strong>Request ID:</strong> %s<br>
			</div>
			<div class="footer">
				Security protection by Gasha WAF &mdash; Request ID: %s
			</div>
		</div>
	</body>
	</html>
	`, RequestID, RequestID)

	_, err := w.Write([]byte(page))
	if err != nil {
		log.Printf("Failed to write custom 403 page: %v", err)
	}
}
