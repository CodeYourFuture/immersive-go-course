package main

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
)

type Case struct {
	Name           string
	In             string
	ExpectedHeader []string
	ExpectedRows   [][]string
}

func TestConsume(t *testing.T) {
	cases := []Case{
		{
			Name:           "Empty",
			In:             "",
			ExpectedHeader: []string{},
			ExpectedRows:   [][]string{},
		},
		{
			Name:           "Simple",
			In:             "name\nTom\nLucy",
			ExpectedHeader: []string{"name"},
			ExpectedRows: [][]string{
				{"Tom"},
				{"Lucy"},
			},
		},
		{
			Name:           "Multiple columns",
			In:             "name,username\nTom,tom\nLucy,lucy",
			ExpectedHeader: []string{"name", "username"},
			ExpectedRows: [][]string{
				{"Tom", "tom"},
				{"Lucy", "lucy"},
			},
		},
	}

	for _, testCase := range cases {
		r := csv.NewReader(strings.NewReader(testCase.In))

		calls := 0
		consume(r, func(s []string) error {
			if len(s) != len(testCase.ExpectedHeader) {
				t.Fatalf("%v: header row: incorrect length, expected %v, got %v", testCase.Name, len(testCase.ExpectedHeader), len(s))
			}
			if !reflect.DeepEqual(s, testCase.ExpectedHeader) {
				t.Fatalf("%v: header row: incorrect value, expected %v, got %v", testCase.Name, testCase.ExpectedHeader, s)
			}
			return nil
		}, func(s []string) error {
			if len(s) != len(testCase.ExpectedRows[calls]) {
				t.Fatalf("%v: row: incorrect length, expected 1, got %v", testCase.Name, len(s))
			}
			if !reflect.DeepEqual(s, testCase.ExpectedRows[calls]) {
				t.Fatalf("%v: row: incorrect value, expected %v, got %v", testCase.Name, testCase.ExpectedRows[calls], s)
			}
			calls += 1
			return nil
		})

		if calls != len(testCase.ExpectedRows) {
			t.Fatalf("%v: incorrect num calls to row processing function, expected %v, got %v", testCase.Name, len(testCase.ExpectedRows), calls)
		}
	}
}
