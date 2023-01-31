package json

import (
	"encoding/json"
	"io"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers"
)

type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]parsers.ScoreRecord, error) {
	var records []parsers.ScoreRecord
	if err := json.NewDecoder(r).Decode(&records); err != nil {
		return nil, err
	}
	return records, nil
}
