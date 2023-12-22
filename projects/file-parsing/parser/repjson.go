package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func ParseRepeatedJson(r io.Reader) ([]Player, error) {
	var players []Player

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}

		var player Player
		err := json.Unmarshal([]byte(scanner.Text()), &player)
		if err != nil {
			return players, fmt.Errorf("json unmarshal error: %w", err)
		}
		players = append(players, player)
	}

	return players, nil
}
