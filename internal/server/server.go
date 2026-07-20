package server

import (
	"net"
	"strconv"
	"sync/atomic"
	"log"
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
	_, err := conn.Write([]byte(`HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13

Hello World!
`))
	if err != nil {
		log.Printf("fail writing response: %s\n", err)
	}
	conn.Close()
}
