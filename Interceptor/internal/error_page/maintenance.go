package error_page

import (
	"fmt"
	"net/http"
)

// SendMaintenanceResponse sends an HTTP 503 response with a maintenance HTML page.
func SendMaintenanceResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
	<title>Under Maintenance</title>
</head>
<body>
	<h1>Site Under Maintenance</h1>
	<p>We're currently performing maintenance. Please check back soon.</p>
</body>
</html>`)
}
