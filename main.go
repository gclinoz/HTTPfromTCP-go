package main

import (
	"strings"
	"errors"
	"fmt"
	"io"
	"os"
	"log"
)

const fileName = "./messages.txt"

func main() {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error when opening %s: %s\n", fileName, err)
	}
	defer f.Close()

	currentLine := ""
	for {
		b := make([]byte, 8, 8)
		n, err := f.Read(b)
		if err != nil {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
			}
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("fail to read content: %s\n", err.Error())
			break
		}
		str := string(b[:n])
		parts := strings.Split(str, "\n")

		for i := 0; i < len(parts) - 1; i++ {
			currentLine += parts[i]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
		}
		currentLine += parts[len(parts) - 1]
	}
}
