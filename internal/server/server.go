package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type state int

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}

		hErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}

	response.WriteStatusLine(conn, response.OK)
	h := response.GetDefaultHeaders(len(buf.Bytes()))
	response.WriteHeaders(conn, h)

	conn.Write(buf.Bytes())
	return
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h *HandlerError) Write(conn io.Writer) {
	response.WriteStatusLine(conn, h.StatusCode)
	head := response.GetDefaultHeaders(len(h.Message))
	response.WriteHeaders(conn, head)
	conn.Write([]byte(h.Message))
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
