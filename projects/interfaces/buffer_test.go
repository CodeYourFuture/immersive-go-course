package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialBytesAreRead(t *testing.T) {
	want := []byte("hello")

	b := NewBufferString("hello")

	got := b.Bytes()

	require.Equal(t, want, got)
}

func TestSubsequentWritesAreAppended(t *testing.T) {
	want := []byte("hello world")

	b := NewBufferString("hello")

	_, err := b.Write([]byte(" world"))
	require.NoError(t, err)

	got := b.Bytes()

	require.Equal(t, want, got)
}

func TestReadWithSliceBigEnoughForWholeBuffer(t *testing.T) {
	b := NewBufferString("hello world")

	slice := make([]byte, 50)

	n, err := b.Read(slice)
	require.NoError(t, err)
	require.Equal(t, 11, n)
	require.Equal(t, []byte("hello world"), slice[:n])
}

func TestReadWithSliceSmallerThanWholeBuffer(t *testing.T) {
	b := NewBufferString("hello world")

	slice := make([]byte, 6)

	n, err := b.Read(slice)
	require.NoError(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, []byte("hello "), slice)

	n, err = b.Read(slice)
	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, []byte("world"), slice[:n])
}
