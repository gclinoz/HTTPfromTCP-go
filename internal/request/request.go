package request

import (
	"errors"
	"bytes"
	"io"
	"fmt"
	"strings"
)

type StateRequest int

const (
	INIT StateRequest = iota
	DONE
)

type Request struct {
	RequestLine RequestLine
	Tracker		StateRequest
}

type RequestLine struct {
	HttpVersion		string
	RequestTarget	string
	Method			string
}

const (
	crlf = "\r\n"
	bufferSize = 8
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0 // keep track of how much data we've read
	req := &Request{ Tracker: INIT }

	for req.Tracker == INIT {
		// if buffer is full, grow it and copy the old data into it
		if readToIndex >= len(buf) {
			newbuf := make([]byte, len(buf) * 2)
			copy(newbuf, buf)
			buf = newbuf
		}

		// read from io into the buffer starting at readToIndex
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.Tracker = DONE
				break
			}
			return nil, err
		}
		readToIndex += n

		numParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numParsed:])
		readToIndex -= numParsed
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.Tracker {
	case INIT:
		out, n, err := parseHTTPbyte(data)
		if err != nil {
			return 0, err
		}
		if n == 0 && err == nil {
			return 0, nil
		}
		r.RequestLine = *out
		r.Tracker = DONE
		return n, nil
	case DONE:
		return 0, fmt.Errorf("trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseHTTPbyte(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	reqText := string(data[:idx])
	reqLine, err := parseRequestLine(reqText)
	if err != nil {
		return nil, 0, err
	}
	return reqLine, idx + 2, nil
}

func parseRequestLine(content string) (*RequestLine, error) {
	if content == "" {
		return nil, fmt.Errorf("empty content")
	}

	// always 3 parts in request line
	partsReq := strings.Split(content, " ")
	if len(partsReq) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", content)
	}

	// check method part only contains capital alphabetic characters
	if ok := isUpper(partsReq[0]); !ok {
		return nil, fmt.Errorf("invalid request method")
	}

	// check http version
	h := strings.Split(partsReq[2], "/")[0]
	if h != "HTTP" {
		return nil, fmt.Errorf("invalid http version")
	}

	// check http version
	v := strings.Split(partsReq[2], "/")[1]
	if v != "1.1" {
		return nil, fmt.Errorf("invalid http version")
	}

	result := &RequestLine{
		HttpVersion:	v,
		RequestTarget:	partsReq[1],
		Method:			partsReq[0],
	}
	return result, nil
}

func isUpper(s string) bool {
	for _, letter := range s {
		if letter < 'A' || letter > 'Z' {
			return false
		}
	}
	return true
}
