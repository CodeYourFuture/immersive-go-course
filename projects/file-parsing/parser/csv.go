package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func ParseCsv(r io.Reader) ([]Player, error) {
	var players []Player

	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return players, fmt.Errorf("error reading csv: %w", err)
	}

	for _, record := range records[1:] {
		highScore, err := strconv.Atoi(record[1])
		if err != nil {
			return players, fmt.Errorf("error parsing a high score: %w", err)
		}
		players = append(players, Player{
			Name:      record[0],
			HighScore: highScore,
		})
	}

	return players, nil
}
