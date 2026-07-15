package main

import (
	"fmt"
	"net"
	"bufio"
	"log"
	"os"
)

const connStr = "localhost:42069"

func main() {
	addr, err := net.ResolveUDPAddr("udp", connStr)
	if err != nil {
		log.Fatalf("fail to resolve address: %s\n", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("fail to create UDP connection: %s\n", err)
	}
	defer conn.Close()

	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("> ")
		ln, err := rd.ReadString('\n')
		if err != nil {
			log.Fatalf("error when reading line: %s\n")
		}
		_, err = conn.Write([]byte(ln))
		if err != nil {
			log.Fatalf("error when writing to UDP: %s\n", err)
		}
		fmt.Printf("Message sent: %s", ln)
	}
}
