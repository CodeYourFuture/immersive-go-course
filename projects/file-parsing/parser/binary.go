package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type byteToUint32Converter interface {
	Uint32(b []byte) uint32
}

func ParseCustomBinaryBigEndian(r io.Reader) ([]Player, error) {
	return ParseCustomBinary(r, []byte{0xFE, 0xFF}, binary.BigEndian)
}

func ParseCustomBinaryLittleEndian(r io.Reader) ([]Player, error) {
	return ParseCustomBinary(r, []byte{0xFF, 0xFE}, binary.LittleEndian)
}

func ParseCustomBinary(r io.Reader, header []byte, converter byteToUint32Converter) ([]Player, error) {
	var players []Player

	data, err := io.ReadAll(r)
	if err != nil {
		return players, fmt.Errorf("error reading: %w", err)
	}

	if !bytes.Equal(data[:2], header) {
		return players, fmt.Errorf("unknown binary format")
	}

	for offset := 2; offset < len(data); {
		highScore := converter.Uint32(data[offset : offset+4])
		offset += 4

		nameEndIndex := findStringEndIndex(data, offset)
		name := string(data[offset:nameEndIndex])
		offset = nameEndIndex + 1

		players = append(players, Player{
			Name:      name,
			HighScore: int(int32(highScore)),
		})
	}

	return players, nil
}

func findStringEndIndex(data []byte, offset int) int {
	for end := offset; end < len(data); end++ {
		if data[end] == 0x00 {
			return end
		}
	}
	return -1
}
