package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteNoNumbers(t *testing.T) {
	buf := bytes.NewBufferString("")

	fp := NewFilteringPipe(buf)

	n, err := fp.Write([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, "hello", buf.String())
}

func TestWriteJustNumbers(t *testing.T) {
	buf := bytes.NewBufferString("")

	fp := NewFilteringPipe(buf)

	n, err := fp.Write([]byte("123"))
	require.NoError(t, err)
	require.Equal(t, 3, n)
	require.Equal(t, "", buf.String())
}

func TestMultipleWrites(t *testing.T) {
	buf := bytes.NewBufferString("")

	fp := NewFilteringPipe(buf)

	n, err := fp.Write([]byte("start="))
	require.NoError(t, err)
	require.Equal(t, 6, n)
	n, err = fp.Write([]byte("1, end=10"))
	require.NoError(t, err)
	require.Equal(t, 9, n)
	require.Equal(t, "start=, end=", buf.String())
}

func TestWriteMixedNumbersAndLetters(t *testing.T) {
	buf := bytes.NewBufferString("")

	fp := NewFilteringPipe(buf)

	n, err := fp.Write([]byte("start=1, end=10"))
	require.NoError(t, err)
	require.Equal(t, 15, n)
	require.Equal(t, "start=, end=", buf.String())
}
