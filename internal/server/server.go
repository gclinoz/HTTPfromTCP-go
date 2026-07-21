package server

import (
	"net"
	"sync/atomic"
	"log"
	"bytes"
	"fmt"

	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

type Server struct {
	listener	net.Listener
	tracker		atomic.Bool
	handler		Handler
}

func Serve(port int, hand Handler) (*Server, error) {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener:	list,
		handler:	hand,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.tracker.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.tracker.Load() {
				return
			}
			log.Printf("fail waiting for connection: %s\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("fail reading request\n")
	}

	var b bytes.Buffer
	herr := s.handler(&b, req)
	if herr != nil {
		err := writeError(conn, herr)
		if err != nil {
			log.Printf("fail to write request error\n")
		}
		return
	}

	h := response.GetDefaultHeaders(b.Len())

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("fail writing status line: %s\n", err)
	}
	err = response.WriteHeaders(conn, h)
	if err != nil {
		log.Printf("fail writing headers: %s\n", err)
	}
	err = response.WriteBody(conn, b.Bytes())
	if err != nil {
		log.Printf("fail writing body: %s\n", err)
	}
}
