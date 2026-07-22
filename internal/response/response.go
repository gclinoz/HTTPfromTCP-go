package response

import (
	"io"
	"fmt"
	"strconv"

	"github.com/gclinoz/HTTPfromTCP-go/internal/headers"
)

type StateCode int

const (
	StatusOK	StateCode = 200
	StatusBad	StateCode = 400
	StatusError	StateCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StateCode) error {
	rp := ""
	switch statusCode {
	case StatusOK:
		rp = "HTTP/1.1 200 OK\r\n"
	case StatusBad:
		rp = "HTTP/1.1 400 Bad Request\r\n"
	case StatusError:
		rp = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		rp = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}
	_, err := w.Write([]byte(rp))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	h := ""
	for k, v := range headers {
		h += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	h += "\r\n"

	_, err := w.Write([]byte(h))
	if err != nil {
		return err
	}
	return nil
}

func WriteResp(w io.Writer, statusCode StateCode, headers headers.Headers, body []byte) error {
	err := WriteStatusLine(w, statusCode)
	if err != nil {
		return err
	}
	err = WriteHeaders(w, headers)
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	err = WriteHeaders(w, headers)
	if err != nil {
		return err
	}
	return nil
}
