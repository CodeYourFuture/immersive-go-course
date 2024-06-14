package binary

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers"
)

type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]parsers.ScoreRecord, error) {
	bufRead := bufio.NewReader(r)

	var records []parsers.ScoreRecord

	byteOrder, err := parseByteOrder(r)
	if err != nil {
		return nil, fmt.Errorf("failed to determine endian-ness: %w", err)
	}

	for {
		if _, err := bufRead.Peek(1); errors.Is(err, io.EOF) {
			break
		}
		var score int32
		if err := binary.Read(bufRead, byteOrder, &score); err != nil {
			return nil, fmt.Errorf("failed to parse score: %w", err)
		}
		nameWithTrailingNull, err := bufRead.ReadString('\x00')
		if err != nil {
			return nil, fmt.Errorf("failed to parse name: %w", err)
		}
		name := nameWithTrailingNull[:len(nameWithTrailingNull)-1]
		records = append(records, parsers.ScoreRecord{
			Name:      name,
			HighScore: score,
		})
	}

	return records, nil
}

func parseByteOrder(r io.Reader) (binary.ByteOrder, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	if buf[0] == '\xFE' && buf[1] == '\xFF' {
		return binary.BigEndian, nil
	} else if buf[0] == '\xFF' && buf[1] == '\xFE' {
		return binary.LittleEndian, nil
	} else {
		return nil, fmt.Errorf("didn't recognise byte-order mark")
	}
}
