package utils

import "net/http"

func GetRequestBodySizeMB(r *http.Request) float64 {
	contentLength := r.ContentLength // size in bytes
	if contentLength <= 0 {
		return 0
	}
	return float64(contentLength) / (1024 * 1024)
}
