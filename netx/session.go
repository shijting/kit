package netx

import (
	"container/list"
	"github.com/shijting/kit/codex"
	"net"
	"sync"
	"sync/atomic"
)

type Session struct {
	net.Conn

	id       uint64
	codex    codex.Codex
	sendChan chan any
	recvMu   sync.Mutex
	sendMu   sync.RWMutex

	closeFlag      int32
	closeChan      chan int
	closeMutex     sync.Mutex
	closeCallbacks *list.List

	State atomic.Value
}
