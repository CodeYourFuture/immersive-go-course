package repeated_json

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers"
)

type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]parsers.ScoreRecord, error) {
	var records []parsers.ScoreRecord
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		var record parsers.ScoreRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
