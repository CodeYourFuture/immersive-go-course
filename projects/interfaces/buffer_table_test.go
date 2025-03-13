package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOurByteBuffer_Table(t *testing.T) {
	for name, tc := range map[string]struct {
		initialContent string
		operations     []operation
	}{
		"read_initial_bytes": {
			initialContent: "hello",
			operations: []operation{
				&bytesOperation{
					wantValue: "hello",
				},
			},
		},
		"subsequent_writes_are_appended": {
			initialContent: "hello",
			operations: []operation{
				&writeOperation{
					value: " world",
				},
				&bytesOperation{
					wantValue: "hello world",
				},
			},
		},
		"read_oversized_slice": {
			initialContent: "hello world",
			operations: []operation{
				&readOperation{
					bufferSize: 50,
					wantValue:  "hello world",
				},
			},
		},
		"read_undersized_slices": {
			initialContent: "hello world",
			operations: []operation{
				&readOperation{
					bufferSize: 6,
					wantValue:  "hello ",
				},
				&readOperation{
					bufferSize: 5,
					wantValue:  "world",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			b := NewBufferString(tc.initialContent)
			for _, operation := range tc.operations {
				operation.Do(t, &b)
			}
		})
	}
}

type operation interface {
	Do(t *testing.T, b *OurByteBuffer)
}

type readOperation struct {
	bufferSize int
	wantValue  string
}

func (o *readOperation) Do(t *testing.T, b *OurByteBuffer) {
	t.Helper()
	buf := make([]byte, o.bufferSize)
	n, err := b.Read(buf)
	require.NoError(t, err)
	require.Equal(t, len(o.wantValue), n)
	require.Equal(t, o.wantValue, string(buf[:n]))
}

type writeOperation struct {
	value string
}

func (o *writeOperation) Do(t *testing.T, b *OurByteBuffer) {
	t.Helper()
	n, err := b.Write([]byte(o.value))
	require.NoError(t, err)
	require.Equal(t, len(o.value), n)
}

type bytesOperation struct {
	wantValue string
}

func (o *bytesOperation) Do(t *testing.T, b *OurByteBuffer) {
	t.Helper()
	got := string(b.Bytes())
	require.Equal(t, o.wantValue, got)
}
