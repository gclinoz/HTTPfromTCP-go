package headers

import (
	"fmt"
	"bytes"
	"strings"
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

	headText := string(data[:idx])
	// CRLF found at the start of the data
	if headText == "" {
		return idx + 2, true, nil
	}

	// parse field-line one key-value pair at a time
	parts := strings.Split(headText, ":")
	if len(parts) == 1 {
		return 0, false, fmt.Errorf("invalid field-line format")
	}
	invalHead := strings.HasPrefix(parts[0], " ")
	invalTail := strings.HasSuffix(parts[0], " ")
	invalIn := strings.Contains(parts[0], " ")
	if invalHead || invalTail || invalIn {
		return 0, false, fmt.Errorf("invalid white space of field-name")
	}
	key := parts[0]
	val := strings.TrimSpace(parts[1])
	if len(parts) > 2 {
		for _, item := range parts[2:] {
			val += ":" + strings.TrimSpace(item)
		}
	}
	h[key] = val
	return idx + 2, false, nil
}
