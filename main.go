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

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		currentLine := ""
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("fail to read content: %s\n", err.Error())
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
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error when opening %s: %s\n", fileName, err)
	}

	ch := getLinesChannel(f)
	for item := range ch {
		fmt.Printf("read: %s\n", item)
	}
}
