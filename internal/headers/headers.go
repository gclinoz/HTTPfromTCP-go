package headers

import (
	"fmt"
	"bytes"
	"strings"
	"slices"
	"log"
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
	existVal, ok := h[key]
	if ok {
		value = existVal + "," + value
	}
	h[key] = value
}

func (h Headers) Del(key string) {
	delete(h, key)
}

func (h Headers) Replace(key, value string) {
	key = strings.ToLower(key)
	_, ok := h[key]
	if !ok {
		log.Printf("replacing a non-exist key\n")
	}
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	val, ok := h[key]
	return val, ok
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
	c >= '0' && c <='9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}
