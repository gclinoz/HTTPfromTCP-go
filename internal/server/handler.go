package server

import (
	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
type ErrorHandler func(w *response.Writer)
