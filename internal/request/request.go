package request

import (
	"errors"
	"bytes"
	"io"
	"fmt"
	"strings"
	"strconv"

	"github.com/gclinoz/HTTPfromTCP-go/internal/headers"
)

type StateRequest int

const (
	requestStateInit StateRequest = iota
	requestStateDone
	requestStateHead
	requestStateBody
)

type Request struct {
	RequestLine RequestLine
	Headers		headers.Headers
	Body		[]byte
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
	requiredHeaderKey = "Content-Length"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0 // keep track of how much data we've read
	req := &Request{
		Headers: headers.NewHeaders(), // program panic without initialize the map
		Tracker: requestStateInit,
	}

	// set condition to req.Tracker == requestStateInit is a bug. It will never parse header
	for req.Tracker != requestStateDone {
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
				val := req.Headers.Get(requiredHeaderKey)
				lenHeader, err := strconv.Atoi(val)
				if err != nil {
					return nil, err
				}
				if len(req.Body) < lenHeader {
					return nil, fmt.Errorf("data shorter than content-length header")
				}
				req.Tracker = requestStateDone
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
	totalBytesParsed := 0
	for r.Tracker != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break // break out the parse loop and read more data
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.Tracker {
	case requestStateInit:
		out, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil // need more data
		}
		r.RequestLine = *out
		r.Tracker = requestStateHead
		return n, nil
	case requestStateHead:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.Tracker = requestStateBody
		}
		return n, nil
	case requestStateBody:
		val := r.Headers.Get(requiredHeaderKey)
		if val == "" {
			r.Tracker = requestStateDone
			return 0, nil
		}

		r.Body = append(r.Body, data...)

		lenHeader, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		if len(r.Body) > lenHeader {
			return 0, fmt.Errorf("data longer than content-legnth header")
		}
		if len(r.Body) == lenHeader {
			r.Tracker = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, fmt.Errorf("trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	reqText := string(data[:idx])
	reqLine, err := parseRequestLineString(reqText)
	if err != nil {
		return nil, 0, err
	}
	return reqLine, idx + 2, nil
}

func parseRequestLineString(content string) (*RequestLine, error) {
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
