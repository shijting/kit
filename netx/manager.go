package netx

import (
	"github.com/shijting/kit/codex"
	"github.com/shijting/kit/option"
	"net"
	"sync"
)

const sessionMapSize uint64 = 32

type Manager struct {
	sessionMaps map[uint64]*sessionMap
	closeOnce   sync.Once
}

func NewManager() *Manager {
	m := new(Manager)

	for i := uint64(0); i < sessionMapSize; i++ {
		m.sessionMaps[i] = &sessionMap{sessions: make(map[uint64]*Session)}
	}
	return m
}

type sessionMap struct {
	sessions map[uint64]*Session
	sync.RWMutex
	isClosed bool
}

func (m *Manager) NewSession(conn net.Conn, code codex.Codex, sendSize int) *Session {
	opts := make([]option.Option[Session], 0)
	if sendSize > 0 {
		opts = append(opts, WithSendSize(sendSize))
	}
	sess := newSession(code, conn, opts...)

	m.putSession(sess)
	return sess
}

func (m *Manager) putSession(sess *Session) {
	sessMap := m.sessionMaps[sess.id%sessionMapSize]
	sessMap.Lock()
	defer sessMap.Unlock()
	sessMap.sessions[sess.id] = sess
	return
}

func (m *Manager) GetSession(sessionId uint64) *Session {
	sessMap := m.sessionMaps[sessionId%sessionMapSize]
	sessMap.RLock()
	defer sessMap.RUnlock()
	return sessMap.sessions[sessionId]
}

func (m *Manager) DelSession(session *Session) {
	sessMap := m.sessionMaps[session.id%sessionMapSize]
	sessMap.Lock()
	defer sessMap.Unlock()
	//delete(sessMap.sessions, session.id)
	m.delSessionById(sessMap, session.id)
}

func (m *Manager) delSessionById(sessMap *sessionMap, sessionId uint64) {
	delete(sessMap.sessions, sessionId)

}

func (m *Manager) Close() {
	m.closeOnce.Do(func() {
		for _, sessMap := range m.sessionMaps {
			sessMap.Lock()
			sessMap.isClosed = true
			for _, sess := range sessMap.sessions {
				_ = sess.Close()
				m.delSessionById(sessMap, sess.id)
			}
			sessMap.Unlock()
		}
	})
}
