package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gclinoz/HTTPfromTCP-go/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
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
