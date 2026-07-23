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

type writerState int

const (
	stateStatusLine writerState = iota
	stateHeaders
	stateBody
)

type Writer struct {
	pen			io.Writer	
	writerState	writerState
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		pen: conn,
	}
}

func (w *Writer) WriteAll(statusCode StateCode, headers headers.Headers, body []byte) (int, error) {
	err := w.WriteStatusLine(statusCode)
	if err != nil {
		return 0, err
	}
	err = w.WriteHeaders(headers)
	if err != nil {
		return 0, err
	}
	n, err := w.WriteBody(body)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteStatusLine(statusCode StateCode) error {
	if w.writerState != stateStatusLine {
		return fmt.Errorf("you are not expecting to write status line")
	}
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
	_, err := w.pen.Write([]byte(rp))
	if err != nil {
		return err
	}
	w.writerState = stateHeaders
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != stateHeaders {
		return fmt.Errorf("you are not expecting to write headers")
	}

	h := ""
	for k, v := range headers {
		h += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	h += "\r\n"

	_, err := w.pen.Write([]byte(h))
	if err != nil {
		return err
	}
	w.writerState = stateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != stateBody {
		return 0, fmt.Errorf("you are not expecting to write body")
	}
	n, err := w.pen.Write(p)
	if err != nil {
		return 0, err
	}
	w.writerState = stateStatusLine
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != stateBody {
		return 0, fmt.Errorf("you are not expecting to write body")
	}

	size := len(p)
	chunk := fmt.Sprintf("%x\r\n", size)
	n1, err := w.pen.Write([]byte(chunk)) // write size in hex followed by carriage return
	if err != nil {
		return 0, err
	}
	n2, err := w.pen.Write(p) // write chunk data
	if err != nil {
		return 0, err
	}
	n3, err := w.pen.Write([]byte("\r\n")) // write trailing carriage return
	if err != nil {
		return 0, err
	}
	return n1 + n2 + n3, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != stateBody {
		return 0, fmt.Errorf("you are not expecting to write body")
	}

	n, err := w.pen.Write([]byte("0\r\n\r\n")) // write end signal
	if err != nil {
		return 0, err
	}
	w.writerState = stateStatusLine
	return n, nil
}
