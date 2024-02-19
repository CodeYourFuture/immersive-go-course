package binary

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseByteOrderBE(t *testing.T) {
	buf := bytes.NewBuffer([]byte("\xFE\xFF"))
	bo, err := parseByteOrder(buf)
	require.NoError(t, err)
	require.Equal(t, binary.BigEndian, bo)
}

func TestParseByteOrderLE(t *testing.T) {
	buf := bytes.NewBuffer([]byte("\xFF\xFE"))
	bo, err := parseByteOrder(buf)
	require.NoError(t, err)
	require.Equal(t, binary.LittleEndian, bo)
}

func TestParseByteOrderWrong(t *testing.T) {
	buf := bytes.NewBuffer([]byte("\xFF\xFF"))
	_, err := parseByteOrder(buf)
	require.ErrorContains(t, err, "didn't recognise byte-order mark")
}

func TestParseByteOrderNotEnoughBytes(t *testing.T) {
	buf := bytes.NewBuffer([]byte("\xFF"))
	_, err := parseByteOrder(buf)
	require.ErrorContains(t, err, "EOF")
}

func TestNotEnoughBytesForScore(t *testing.T) {
	buf := bytes.NewBuffer([]byte("\xFE\xFF\x01\x00"))
	parser := &Parser{}
	_, err := parser.Parse(buf)
	require.ErrorContains(t, err, "failed to parse score")
}

func TestMissingNullTerminator(t *testing.T) {
	buf := bytes.NewBuffer([]byte("\xFE\xFF\x01\x00\x00\x00Aya"))
	parser := &Parser{}
	_, err := parser.Parse(buf)
	require.ErrorContains(t, err, "failed to parse name")
}
