package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

type serverState int

const (
	initialized serverState = iota
	running
	down
)

type Server struct {
	State serverState
	Listener net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		State:   initialized,
		Listener: l,
	}

	go s.listen()
	return s, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.State == down {
				return
			}
			log.Printf("accept error: %v", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			io.Copy(c, c)
		}(conn)
	}
}

func (s *Server) Close() error {
	s.State = down
	return s.Listener.Close()
}