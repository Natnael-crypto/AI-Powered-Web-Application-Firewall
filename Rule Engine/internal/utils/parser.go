package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// ParseRequest formats an HTTP request into a presentable string.
func ParseRequest(r *http.Request) string {
	var requestDetails bytes.Buffer

	// Start with the request line
	requestDetails.WriteString(fmt.Sprintf("Method: %s\nURL: %s\nProtocol: %s\n", r.Method, r.URL.String(), r.Proto))

	// Add headers
	requestDetails.WriteString("Headers:\n")
	for name, values := range r.Header {
		for _, value := range values {
			requestDetails.WriteString(fmt.Sprintf("  %s: %s\n", name, value))
		}
	}

	// Add query parameters
	if len(r.URL.Query()) > 0 {
		requestDetails.WriteString("Query Parameters:\n")
		for name, values := range r.URL.Query() {
			for _, value := range values {
				requestDetails.WriteString(fmt.Sprintf("  %s: %s\n", name, value))
			}
		}
	}

	// Add the body (if available)
	if r.Body != nil {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reassign body so it can be read again later
		requestDetails.WriteString("Body:\n")
		requestDetails.WriteString(fmt.Sprintf("%s\n", string(body)))
	}

	return requestDetails.String()
}
