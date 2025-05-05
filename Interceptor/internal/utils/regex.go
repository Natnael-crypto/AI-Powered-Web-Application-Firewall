package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseHeaders(input string) string {
	re := regexp.MustCompile(`(\w[\w-]*):\[(.*?)\]`)

	matches := re.FindAllStringSubmatch(input, -1)

	var headers []string
	for _, match := range matches {
		header := fmt.Sprintf("%s: %s", match[1], match[2])
		headers = append(headers, header)
	}

	return strings.Join(headers, ", ")
}
