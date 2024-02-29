package netx

import (
	"github.com/shijting/kit"
	"github.com/shijting/kit/codex"
	"io"
	"net"
	"strings"
	"time"
)

type Server struct {
	net.Listener
	//manager      *Manager
	codex codex.Codex
	//handler      Handler
	sendChanSize int
}

func NewServer(listener net.Listener, sendChanSize int) *Server {
	return &Server{Listener: listener, sendChanSize: sendChanSize}
}

func (s *Server) Serve() error {
	for {
		conn, err := s.accept()

		if err != nil {
			return err
		}

		_ = conn
	}
}
func (s *Server) accept() (net.Conn, error) {
	retry := kit.NewExponentialBackoffRetry(5*time.Millisecond, time.Second, 5)
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				tempDelay, ok := retry.Next()
				if !ok {
					return nil, err
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}

		return conn, nil
	}
}
