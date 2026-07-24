package main

import (
	"log"
	"strings"
	"net/http"
	"errors"
	"io"

	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

const bufSize = 1024

func handlerMain(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
		handlerErrorRequest(w)
		return
	}

	if target == "/myproblem" {
		handlerErrorInternal(w)
		return
	}

	if strings.HasPrefix(target, "/httpbin/stream") {
		path:= strings.TrimPrefix(target, "/httpbin")
		handlerProxyStream(w, path)
		return
	}

	if strings.HasPrefix(target, "/httpbin/html") {
		handlerProxy(w)
		return
	}

	h := response.GetDefaultHeaders(len(pass))
	h.Replace("Content-Type", "text/html")
	n, err := w.WriteAll(response.StatusOK, h, []byte(pass))
	log.Printf("wrote %d of %d bytes:", n, len(pass))
	if err != nil {
		log.Printf("fail to write response: %s\n", err)
	}
}

func handlerErrorRequest(w *response.Writer) {
	h := response.GetDefaultHeaders(len(bad))
	h.Replace("Content-Type", "text/html")
	n, err := w.WriteAll(response.StatusBad, h, []byte(bad))
	log.Printf("wrote %d of %d bytes:", n, len(bad))
	if err != nil {
		log.Printf("fail to write response: %s\n", err)
	}
}

func handlerErrorInternal(w *response.Writer) {
	h := response.GetDefaultHeaders(len(internal))
	h.Replace("Content-Type", "text/html")
	n, err := w.WriteAll(response.StatusError, h, []byte(internal))
	log.Printf("wrote %d of %d bytes:", n, len(internal))
	if err != nil {
		log.Printf("fail to write response: %s\n", err)
	}
}

func handlerProxyStream(w *response.Writer, path string) {
	addr := "https://httpbin.org"
	URL := addr + path
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("fail to get httpbin content: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("httpbin returned status %d\n", resp.StatusCode)
		w.WriteStatusLine(response.StatusError)
		w.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}

	h := response.GetDefaultHeaders(0)
	h.Del("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	err = w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Printf("fail to write status line\n")
		return
	}
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("fail to write headers\n")
		return
	}

	buf := make([]byte, bufSize)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Printf("error when writing chunk: %s\n", err)
				return
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("error when reading httpbin response: %s\n", err)
			return
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("fail to terminate chunk write: %s\n", err)
	}
	err = w.WriteTrailers(h)
	if err != nil {
		log.Printf("fail to write trailers: %s\n", err)
	}
}

func handlerProxy(w *response.Writer) {
	addr := "https://httpbin.org"
	URL := addr + "/html"
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("fail to get httpbin content: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("httpbin returned status %d\n", resp.StatusCode)
		w.WriteStatusLine(response.StatusError)
		w.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}

	h := response.GetDefaultHeaders(0)
	h.Del("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")
	err = w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Printf("fail to write status line\n")
		return
	}
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("fail to write headers\n")
		return
	}

	buf := make([]byte, bufSize)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Printf("error when writing chunk: %s\n", err)
				return
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("error when reading httpbin response: %s\n", err)
			return
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("fail to terminate chunk write: %s\n", err)
	}
	err = w.WriteTrailers(h)
	if err != nil {
		log.Printf("fail to write trailers: %s\n", err)
	}
}
