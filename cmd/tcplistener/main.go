package main

import (
	"net"
	"strings"
	"errors"
	"fmt"
	"io"
	"log"
)

const addr = "127.0.0.1:42069"

func getLinesChannel(c net.Conn) <-chan string {
	ch := make(chan string)

	go func() {
		defer c.Close()
		defer close(ch)

		currentLine := ""
		for {
			b := make([]byte, 8)
			n, err := c.Read(b)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("fail to read content: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts) - 1; i++ {
				ch <- currentLine + parts[i]
				currentLine = ""
			}
			currentLine += parts[len(parts) - 1]
		}
	}()
	return ch
}

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

		ch := getLinesChannel(conn)
		for item := range ch {
			fmt.Println(item)
		}
		fmt.Println("Connection closed...")
	}
}
