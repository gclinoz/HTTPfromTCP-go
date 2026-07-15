package request

import (
	"io"
	"fmt"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion		string
	RequestTarget	string
	Method			string
}

func isUpper(s string) bool {
	for _, letter := range s {
		if !unicode.IsUpper(letter) && unicode.IsLetter(letter) {
			return false
		}
	}
	return true
}

func parseRequestLine(content string) (RequestLine, error) {
	if content == "" {
		return RequestLine{}, fmt.Errorf("empty content")
	}

	parts := strings.Split(content, "\r\n")
	reqLine := parts[0]

	// always 3 parts in request line
	partsReq := strings.Split(reqLine, " ")
	if len(partsReq) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line")
	}

	// check method part only contains capital alphabetic characters
	if ok := isUpper(partsReq[0]); !ok {
		return RequestLine{}, fmt.Errorf("invalid request method")
	}

	// check http version part is 1.1
	v := strings.Split(partsReq[2], "/")[1]
	if v != "1.1" {
		return RequestLine{}, fmt.Errorf("invalid http version")
	}

	result := RequestLine{
		HttpVersion:	v,
		RequestTarget:	partsReq[1],
		Method:			partsReq[0],
	}
	return result, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}
	contentStr := string(b)

	out, err := parseRequestLine(contentStr)
	if err != nil {
		return &Request{}, err
	}
	final := &Request{ RequestLine: out }
	return final, nil
}
