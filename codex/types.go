package codex

type Codex interface {
	Send(any) error
	Receive(any) error
	Close() error
}
