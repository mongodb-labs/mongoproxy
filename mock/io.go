// Package mock contains mock structs that implement interfaces for reading /
// writing data, which is useful for testing purposes.
package mock

import (
	"io"
)

// A MockIO emulates an io.Reader and an io.Writer for testing purposes
// the bytes of a MockIO can be immediately examined, and also prefilled
// ahead of time.
type MockIO struct {
	Input  []byte
	Output []byte
	pos    int
}

func (m *MockIO) New() {
	m.pos = 0
}

func (m *MockIO) Read(b []byte) (int, error) {
	inputLen := len(m.Input) - m.pos
	bLen := len(b)
	n := 0

	for i := 0; i < inputLen && i < bLen; i++ {
		b[i] = m.Input[m.pos]

		m.pos += 1
		n += 1
	}

	return n, io.EOF
}
func (m *MockIO) Write(b []byte) (int, error) {
	m.Output = append(m.Output, b...)
	return len(b), nil
}
