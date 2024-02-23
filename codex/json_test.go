package codex

import (
	"bytes"
	"errors"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestJson(t *testing.T) {
	type Message1 struct {
		Field1 string
		Field2 int
		Field3 []byte
	}
	type message2 struct {
		Field1 string
		Field2 int
		Field3 []byte
	}
	var str string
	testCases := []struct {
		name    string
		input   any
		want    any
		wantErr error
	}{
		{
			name:    "nil input",
			input:   nil,
			want:    nil,
			wantErr: errors.New("input nil value"),
		},
		{
			name:    "empty input",
			input:   &Message1{},
			want:    &Message1{},
			wantErr: nil,
		},
		{
			name: "normal input",
			input: &Message1{
				Field1: "abc",
				Field2: 123,
				Field3: []byte("hello"),
			},
			want: &Message1{
				Field1: "abc",
				Field2: 123,
				Field3: []byte("hello"),
			},
			wantErr: nil,
		},
		{
			name:    "string input",
			input:   "hello",
			want:    &str,
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			protocol := NewJson(buf)
			err := protocol.Send(tc.input)
			assert.Equal(t, err, tc.wantErr)
			err = protocol.Receive(tc.want)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, tc.input, tc.want)
		})
	}
}
