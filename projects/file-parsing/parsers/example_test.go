package parsers_test

import (
	"os"
	"testing"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/binary"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/csv"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/json"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/repeated_json"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	wantRecords := []parsers.ScoreRecord{
		{Name: "Aya", HighScore: 10},
		{Name: "Prisha", HighScore: 30},
		{Name: "Charlie", HighScore: -1},
		{Name: "Margot", HighScore: 25},
	}

	for name, tc := range map[string]struct {
		filename string
		parser   parsers.Parser
	}{
		"binary-be": {
			filename: "custom-binary-be.bin",
			parser:   &binary.Parser{},
		},
		"binary-le": {
			filename: "custom-binary-le.bin",
			parser:   &binary.Parser{},
		},
		"csv": {
			filename: "data.csv",
			parser:   &csv.Parser{},
		},
		"json": {
			filename: "json.txt",
			parser:   &json.Parser{},
		},
		"repeated_json": {
			filename: "repeated-json.txt",
			parser:   &repeated_json.Parser{},
		},
	} {
		f, err := os.Open("../examples/" + tc.filename)
		require.NoError(t, err, "%s: failed to open %s", name, tc.filename)
		gotRecords, err := tc.parser.Parse(f)
		require.NoError(t, err, "%s: failed to parse %s", name, tc.filename)
		require.Equal(t, wantRecords, gotRecords, "%s: wrong records returned", name)
	}
}
