package server

import (
	"io"
	"log"

	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

type HandlerError struct {
	Status	response.StateCode
	Message	string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) Write(w io.Writer) {
	h := response.GetDefaultHeaders(len(he.Message))
	b := []byte(he.Message)
	err := response.WriteResp(w, he.Status, h, b)
	if err != nil {
		log.Println("Error when writing response")
	}
}
