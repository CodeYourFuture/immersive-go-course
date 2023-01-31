package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers"
)

type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]parsers.ScoreRecord, error) {
	reader := csv.NewReader(r)
	headerRecord, err := reader.Read()
	if err != nil {
		return nil, err
	}
	nameIndex := -1
	highScoreIndex := -1
	for i, col := range headerRecord {
		if col == "name" {
			nameIndex = i
		} else if col == "high score" {
			highScoreIndex = i
		} else {
			return nil, fmt.Errorf("unexpected header %q - expected %q and %q", col, "name", "high score")
		}
	}
	if nameIndex == -1 || highScoreIndex == -1 {
		return nil, fmt.Errorf("incorrect headers - expected to find %q and %q", "name", "high score")
	}

	var records []parsers.ScoreRecord
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, err
			}
		}
		highScoreString := record[highScoreIndex]
		highScore, err := strconv.ParseInt(highScoreString, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("saw high score which wasn't an int32: %q: %w", highScoreString, err)
		}
		records = append(records, parsers.ScoreRecord{
			Name:      record[nameIndex],
			HighScore: int32(highScore),
		})
	}
	return records, nil
}

// TODO: Add some edge-case tests
