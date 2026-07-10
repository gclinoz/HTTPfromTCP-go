package main

import (
	"fmt"
	"io"
	"os"
	"log"
)

const fileName = "./messages.txt"

func main() {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error when opening the file")
	}

	container := make([]byte, 8)
	for {
		_, err = f.Read(container)
		if err == io.EOF {
			return
		}
		fmt.Printf("read: %s\n", container)
	}
}
