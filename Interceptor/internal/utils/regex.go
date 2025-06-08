package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseHeaders(input string) string {
	re := regexp.MustCompile(`(\w[\w-]*):\[(.*?)\]`)

	skipHeaders := map[string]bool{
		"Sec-Ch-Ua":                 true,
		"Sec-Ch-Ua-Mobile":          true,
		"Sec-Ch-Ua-Platform":        true,
		"Sec-Fetch-Dest":            true,
		"Sec-Fetch-Mode":            true,
		"Sec-Fetch-Site":            true,
		"Sec-Fetch-User":            true,
		"Upgrade-Insecure-Requests": true,
		"Dnt":                       true,
		"Te":                        true,
	}

	matches := re.FindAllStringSubmatch(input, -1)

	var headers []string
	for _, match := range matches {
		headerName := match[1]
		headerValue := match[2]

		if skipHeaders[headerName] {
			continue
		}

		header := fmt.Sprintf("%s: %s", headerName, headerValue)
		headers = append(headers, header)
	}

	return strings.Join(headers, ", ")
}
