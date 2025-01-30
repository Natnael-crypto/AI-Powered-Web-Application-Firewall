package parser

import (
	"fmt"
	"io"
	"net/http"
)

// ParseRequest extracts metadata from an HTTP request
func ParseRequest(r *http.Request) map[string]string {
	metadata := make(map[string]string)
	metadata["Method"] = r.Method
	metadata["URL"] = r.URL.String()
	metadata["Headers"] = fmt.Sprintf("%v", r.Header)
	metadata["Body"] = parseRequestBody(r)
	return metadata
}

// parseRequestBody extracts and parses the request body
func parseRequestBody(r *http.Request) string {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body)
}
