package netx

import (
	"github.com/shijting/kit/codex"
	"github.com/shijting/kit/option"
	"net"
	"time"
)

// Listen listens on the network address addr and then calls Serve with handler to handle requests on incoming connections.
func Listen(addr string, protocol string, code codex.Codex, handler Handler, sendSize int) (*Server, error) {
	listener, err := net.Listen(protocol, addr)
	if err != nil {
		return nil, err
	}

	return NewServer(listener, code, handler, sendSize), nil
}

// Dial connects to the address on the named network.
func Dial(addr string, protocol string, code codex.Codex, sendSize int) (*Session, error) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		return nil, err
	}

	opts := make([]option.Option[Session], 0)
	if sendSize > 0 {
		opts = append(opts, WithSendSize(sendSize))
	}
	return NewSession(code, conn, opts...), nil
}

// DialTimeout connects to the address on the named network with a timeout.
func DialTimeout(addr string, protocol string, code codex.Codex, sendSize int, timeout time.Duration) (*Session, error) {
	conn, err := net.DialTimeout(protocol, addr, timeout)
	if err != nil {
		return nil, err
	}

	opts := make([]option.Option[Session], 0)
	if sendSize > 0 {
		opts = append(opts, WithSendSize(sendSize))
	}
	return NewSession(code, conn, opts...), nil
}
