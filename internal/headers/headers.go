package headers

import (
	"fmt"
	"bytes"
	"strings"
	"slices"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil // no CRLF found, need more data
	}
	if idx == 0 {
		return 2, true, nil // CRLF found at the start, headers are done
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	if len(parts) < 2 {
		return 0, false, fmt.Errorf("invalid header: %s", parts[0]) // no colon found
	}

	key := string(parts[0])
	if strings.ContainsAny(key, " \t") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	key = strings.TrimSpace(key)
	val := bytes.TrimSpace(parts[1])
	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}
	h.Set(key, string(val))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
	c >= 'a' && c <= 'z' ||
	c >= '0' && c <='0' {
		return true
	}

	return slices.Contains(tokenChars, c)
}
