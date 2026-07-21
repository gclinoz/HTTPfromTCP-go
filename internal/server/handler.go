package server

import (
	"io"

	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

type HandlerError struct {
	Status	response.StateCode
	Message	string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func writeError(w io.Writer, herr *HandlerError) error {
	err := response.WriteStatusLine(w, herr.Status)
	if err != nil {
		return err
	}

	h := response.GetDefaultHeaders(len(herr.Message))
	err = response.WriteHeaders(w, h)
	if err != nil {
		return err
	}

	b := []byte(herr.Message)
	err = response.WriteBody(w, b)
	if err != nil {
		return err
	}
	return nil
}
