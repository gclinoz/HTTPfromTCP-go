package server

import (
	"net"
	"strconv"
	"sync/atomic"
	"log"

	"github.com/gclinoz/HTTPfromTCP-go/internal/response"
)

type Server struct {
	listener	net.Listener
	tracker		atomic.Bool
}

func Serve(port int) (*Server, error) {
	portStr := strconv.Itoa(port)
	listener, err := net.Listen("tcp", ":" + portStr)
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
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
				break
			}
			log.Printf("fail waiting for connection: %s\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("fail writing status line: %s\n", err)
	}

	h := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, h)
	if err != nil {
		log.Printf("fail writing headers: %s\n", err)
	}

	conn.Close()
}
