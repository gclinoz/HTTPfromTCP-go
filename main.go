package main

import (
	"fmt"
	"io"
	"os"
	"log"
)

func main() {
	f, err := os.Open("./messages.txt")
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
