package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

func ParseJson(r io.Reader) ([]Player, error) {
	var players []Player

	data, err := io.ReadAll(r)
	if err != nil {
		return players, fmt.Errorf("error reading: %w", err)
	}

	err = json.Unmarshal(data, &players)
	if err != nil {
		return players, fmt.Errorf("json unmarshal error: %w", err)
	}

	return players, nil
}
