package encoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []byte
	}{
		{"normal", "Hello", []byte("Hello")},
		{"empty", "", []byte("")},
		{"chinese", "你好", []byte("你好")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, StringToBytes(tc.input))
		})
	}
}
