package parsers

import (
	"io"
)

type ScoreRecord struct {
	Name      string `json:"name"`
	HighScore int32  `json:"high_score"`
}

type Parser interface {
	Parse(file io.Reader) ([]ScoreRecord, error)
}
