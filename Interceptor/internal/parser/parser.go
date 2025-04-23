package parser

import (
	"fmt"
	"io"
	"net/http"
)

func ParseRequest(r *http.Request) map[string]string {
	metadata := make(map[string]string)
	metadata["Method"] = r.Method
	metadata["URL"] = r.URL.String()
	metadata["Headers"] = fmt.Sprintf("%v", r.Header)
	metadata["Body"] = parseRequestBody(r)
	return metadata
}

func parseRequestBody(r *http.Request) string {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body)
}
