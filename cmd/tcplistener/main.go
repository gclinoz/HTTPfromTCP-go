package main

import (
	"net"
	"fmt"
	"log"

	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
)

const addr = "127.0.0.1:42069"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error when opening listener: %s\n", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("fail waiting for connection: %s\n", err)
		}
		fmt.Println("Connection accepted!")

		data, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("fail reading request")
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", data.RequestLine.Method)
		fmt.Println("- Target:", data.RequestLine.RequestTarget)
		fmt.Println("- Version:", data.RequestLine.HttpVersion)

		fmt.Println("Connection closed...")
	}
}
