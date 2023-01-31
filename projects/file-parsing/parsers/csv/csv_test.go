package csv_test

import (
	"strings"
	"testing"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/csv"
	"github.com/stretchr/testify/require"
)

func TestTooManyColumns(t *testing.T) {
	parser := &csv.Parser{}
	_, err := parser.Parse(strings.NewReader("name,high score,opponent\nAya,12,Prisha\n"))
	require.ErrorContains(t, err, "unexpected header \"opponent\"")
}

func TestNotEnoughColumns(t *testing.T) {
	parser := &csv.Parser{}
	_, err := parser.Parse(strings.NewReader("name\nAya\n"))
	require.ErrorContains(t, err, "incorrect headers")
	require.ErrorContains(t, err, "high score")
}

func TestWrongColumns(t *testing.T) {
	parser := &csv.Parser{}
	_, err := parser.Parse(strings.NewReader("name,low score\nAya,12\n"))
	require.ErrorContains(t, err, "unexpected header")
	require.ErrorContains(t, err, "low score")
	require.ErrorContains(t, err, "high score")
}
