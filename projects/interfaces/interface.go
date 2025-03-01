package main

import (
	"bytes"
	"fmt"
)

// my own bytes buffer
type OurByteBuffer struct {
	bytes        []byte
	readPosition int
}

// NewBufferString creates a new buffer with the given string
func NewBufferString(s string) OurByteBuffer {
	return OurByteBuffer{bytes: []byte(s), readPosition: 0}
}

// Bytes returns the bytes in the buffer
func (b *OurByteBuffer) Bytes() []byte {
	return b.bytes
}

// Read reads the next len(p) bytes from the buffer or until the buffer is drained
func (b *OurByteBuffer) Write(bytes []byte) (n int, err error) {
	b.bytes = append(b.bytes, bytes...)
	return len(bytes), nil
}

// Read reads the next len(p) bytes from the buffer or until the buffer is drained
func (b *OurByteBuffer) Read(p []byte) (n int, err error) {
	reaminingBytes := len(b.bytes) - b.readPosition
	bytesToRead := min(reaminingBytes, len(p))
	copy(p, b.bytes[b.readPosition:b.readPosition+bytesToRead])
	b.readPosition += bytesToRead
	return bytesToRead, nil
}

func main() {
	var b bytes.Buffer

	b.Write([]byte("hello"))
	b.Write([]byte("world"))

	got := b.Bytes()
	want := []byte{}

	fmt.Printf("%v, %v\n", got, want)
}
