package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"io"

	"github.com/gclinoz/HTTPfromTCP-go/internal/server"
	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerTest)
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
	// if req.RequestLine.RequestTarget == "/yourproblem" {
	// 	return &server.HandlerError{
	// 		Status:		response.StatusBad,
	// 		Message:	"Your problem is not my problem\n",
	// 	}
	// }
	// if req.RequestLine.RequestTarget == "/myproblem" {
	// 	return &server.HandlerError{
	// 		Status:		response.StatusError,
	// 		Message:	"Woopsie, my bad\n",
	// 	}
	// }
	// _, err := w.Write([]byte("All good, frfr\n"))
	// if err != nil {
	// 	return &server.HandlerError{
	// 		Status:		response.StatusError,
	// 		Message:	"Woopsie, my bad\n",
	// 	}
	// }
}
