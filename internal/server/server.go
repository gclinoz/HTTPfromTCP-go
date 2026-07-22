package server

import (
	"net"
	"sync/atomic"
	"log"
	"fmt"

	"github.com/gclinoz/HTTPfromTCP-go/internal/request"
	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

type Server struct {
	listener	net.Listener
	tracker		atomic.Bool
	handler		Handler
	errorer		ErrorHandler
}

func Serve(port int, hand Handler, erhand ErrorHandler) (*Server, error) {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener:	list,
		handler:	hand,
		errorer:	erhand,
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

	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		s.errorer(w)
		return
	}

	s.handler(w, req)
}
