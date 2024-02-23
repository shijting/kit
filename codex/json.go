package codex

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	_               Codex = (*Json)(nil)
	ErrJsonNilValue       = errors.New("input nil value")
)

type Json struct {
	encoder *json.Encoder
	decoder *json.Decoder
	closer  io.Closer
}

func NewJson(rw io.ReadWriter) Codex {
	j := &Json{
		encoder: json.NewEncoder(rw),
		decoder: json.NewDecoder(rw),
	}
	j.closer, _ = rw.(io.Closer)
	return j
}

func (j *Json) Send(t any) error {
	if t == nil {
		return ErrJsonNilValue
	}
	return j.encoder.Encode(t)
}

func (j *Json) Receive(t any) error {
	if t == nil {
		return ErrJsonNilValue
	}
	return j.decoder.Decode(t)
}

func (j *Json) Close() error {
	if j.closer != nil {
		return j.closer.Close()
	}
	return nil
}
