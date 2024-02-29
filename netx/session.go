package netx

import (
	"errors"
	"github.com/shijting/kit/codex"
	"github.com/shijting/kit/option"
	"net"
	"sync"
	"sync/atomic"
)

var (
	globalSessionId uint64
)

var (
	ErrSessionClosed = errors.New("session closed")
	ErrSendChanFull  = errors.New("send channel full")
)

type CloseHandler func(session *Session)

type SessionManager interface {
	AddSession(session *Session)
	RemoveSession(session *Session)
	GetSession(id uint64) *Session
}

type Session struct {
	net.Conn
	id     uint64
	codex  codex.Codex
	sendCh chan any

	recvMu    sync.Mutex
	sendMu    sync.RWMutex
	closeFlag int32
	closeCh   chan struct{}
	closeMu   sync.Mutex

	State atomic.Value

	//	关闭回调函数
	closeCallbacks []CloseHandler
}

// todo
func NewSession() {

}

func newSession(codex codex.Codex, conn net.Conn, opts ...option.Option[Session]) *Session {
	sess := &Session{
		Conn:           conn,
		codex:          codex,
		closeCh:        make(chan struct{}),
		closeCallbacks: make([]CloseHandler, 0),
		id:             atomic.AddUint64(&globalSessionId, 1),
	}

	option.Options[Session](opts).Apply(sess)

	go sess.sendLoop()
	return sess
}

func WithSendCh(size int) option.Option[Session] {
	return func(s *Session) {
		s.sendCh = make(chan any, size)
	}
}

func (s *Session) ID() uint64 {
	return s.id
}

func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) == 1
}

func (s *Session) Addr() string {
	return s.RemoteAddr().String()
}
func (s *Session) sendLoop() {
	if s.sendCh == nil {
		return
	}
	defer s.Close()

	for {
		select {
		case <-s.closeCh:
			return
		case msg, ok := <-s.sendCh:
			if !ok || (s.codex.Send(msg)) != nil {
				return
			}
		}
	}
}

func (s *Session) Send(msg any) error {
	if s.sendCh != nil {

		s.sendMu.RLock()
		defer s.sendMu.RUnlock()

		if s.IsClosed() {
			return ErrSessionClosed
		}

		select {
		case s.sendCh <- msg:
			return nil
		default:
			//  send chan full
			return ErrSendChanFull
		}
	}

	s.sendMu.Lock()
	defer s.sendMu.Unlock()

	if s.IsClosed() {
		return ErrSessionClosed
	}

	err := s.codex.Send(msg)
	if err != nil {
		s.Close()
	}
	return err
}

func (s *Session) Receive(a any) error {
	s.recvMu.Lock()
	defer s.recvMu.Unlock()
	if s.IsClosed() {
		return ErrSessionClosed
	}
	err := s.codex.Receive(a)
	if err != nil {
		s.Close()
	}
	return err
}

func (s *Session) Close() error {
	if atomic.CompareAndSwapInt32(&s.closeFlag, 0, 1) {
		// 执行关闭回调函数
		for _, callback := range s.closeCallbacks {
			callback(s)
		}
		s.closeFlag = 1
		// TODO 从管理器中移除
		s.Conn.Close()
		close(s.closeCh)
	}
	return nil
}

func (s *Session) AddCloseCallback(callback CloseHandler) {
	s.closeCallbacks = append(s.closeCallbacks, callback)
}
