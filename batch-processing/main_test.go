package main

import (
	"strings"
	"testing"
)

func TestReadValidateCsv(t *testing.T) {
	in := strings.NewReader(`url
http://host/path.jpg`)
	records, err := readAndValidateCsv(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("records incorrect length: expected 2, got %d\n", len(records))
	}
}

func TestReadValidateCsvHeaderValidation(t *testing.T) {
	in := strings.NewReader(`nope
http://host/path.jpg`)
	records, err := readAndValidateCsv(in)
	if err == nil {
		t.Fatalf("expected error: got %v\n", records)
	}
}

func TestReadValidateCsvEmptyCsv(t *testing.T) {
	in := strings.NewReader("")
	records, err := readAndValidateCsv(in)
	if err == nil {
		t.Fatalf("expected error: got %v\n", records)
	}
}

func TestReadValidateCsvEmptyBody(t *testing.T) {
	in := strings.NewReader(`url`)
	records, err := readAndValidateCsv(in)
	if err == nil {
		t.Fatalf("expected error: got %v\n", records)
	}
}
