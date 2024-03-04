package netx

type Handler interface {
	HandleSession(*Session)
}

type HandlerFunc func(*Session)

func (f HandlerFunc) HandleSession(session *Session) {
	f(session)
}
