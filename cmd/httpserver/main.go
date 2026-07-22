package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gclinoz/HTTPfromTCP-go/internal/server"
	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerTest, handlerErrorRequest)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	// common pattern in Go for gracefully shutting down a server
	// because server.Serve returns immediately, if we exit main immediately,
	// the server will just stop.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerTest(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handlerErrorRequest(w)
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		handlerErrorInternal(w)
	}

	h := response.GetDefaultHeaders(len(pass))
	h.Replace("Content-Type", "text/html")
	err := w.WriteAll(response.StatusOK, h, []byte(pass))
	if err != nil {
		log.Printf("fail to write response: %s", err)
	}
}

func handlerErrorRequest(w *response.Writer) {
	h := response.GetDefaultHeaders(len(bad))
	h.Replace("Content-Type", "text/html")
	err := w.WriteAll(response.StatusBad, h, []byte(bad))
	if err != nil {
		log.Printf("fail to write response: %s", err)
	}
}

func handlerErrorInternal(w *response.Writer) {
	h := response.GetDefaultHeaders(len(internal))
	h.Replace("Content-Type", "text/html")
	err := w.WriteAll(response.StatusError, h, []byte(internal))
	if err != nil {
		log.Printf("fail to write response: %s", err)
	}
}
