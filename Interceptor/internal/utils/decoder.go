package utils

import (
	"encoding/base64"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

func RecursiveDecode(input string, maxDepth int) string {
	current := input

	for i := 0; i < maxDepth; i++ {
		decoded := current

		decoded = urlDecode(decoded)
		decoded = htmlDecode(decoded)
		decoded = unicodeDecode(decoded)
		decoded = hexDecode(decoded)
		decoded = base64Decode(decoded)

		decoded = normalizePayload(decoded)

		if decoded == current {
			break
		}
		current = decoded
	}

	return current
}

func urlDecode(s string) string {
	decoded, err := url.QueryUnescape(s)
	if err == nil {
		return decoded
	}
	return s
}

func htmlDecode(s string) string {
	return html.UnescapeString(s)
}

func unicodeDecode(s string) string {
	re := regexp.MustCompile(`(\\u|%u)([0-9a-fA-F]{4})`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		hex := match[len(match)-4:]
		var r rune
		fmt.Sscanf(hex, "%04x", &r)
		return string(r)
	})
}

func hexDecode(s string) string {
	re := regexp.MustCompile(`\\x([0-9a-fA-F]{2})`)
	s = re.ReplaceAllStringFunc(s, func(match string) string {
		var hexByte byte
		fmt.Sscanf(match, `\x%02x`, &hexByte)
		return string(hexByte)
	})

	if len(s)%2 == 0 && isHexString(s) {
		var decoded strings.Builder
		for i := 0; i < len(s); i += 2 {
			var b byte
			fmt.Sscanf(s[i:i+2], "%02x", &b)
			decoded.WriteByte(b)
		}
		return decoded.String()
	}

	return s
}

func base64Decode(s string) string {
	s = strings.TrimSpace(s)
	if len(s) < 8 || len(s)%4 != 0 {
		return s
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err == nil && isMostlyPrintable(decoded) {
		return string(decoded)
	}
	return s
}

func isMostlyPrintable(data []byte) bool {
	printable := 0
	for _, b := range data {
		if b == 9 || b == 10 || b == 13 || (b >= 32 && b <= 126) {
			printable++
		}
	}
	return float64(printable)/float64(len(data)) > 0.8
}

func isHexString(s string) bool {
	re := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	return re.MatchString(s)
}

func normalizePayload(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return s
}
