package main

import (
	"encoding/csv"
	"errors"
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

func TestHeaderError(t *testing.T) {
	r := csv.NewReader(strings.NewReader("a\nb"))

	expected := "inner header error"
	err := consume(r, func(s []string) error {
		return errors.New(expected)
	}, func(s []string) error {
		t.Fatal("row processing function called unexpectedly")
		return nil
	})

	if err == nil {
		t.Fatal("consume did not return an error")
	}

	inner := errors.Unwrap(err)
	if inner == nil {
		t.Fatal("incorrect error: could not unwrap")
	}
	if inner.Error() != expected {
		t.Fatalf("incorrect error: expected %v, got %v", expected, inner.Error())
	}
}

func TestRowError(t *testing.T) {
	r := csv.NewReader(strings.NewReader("a\nb"))

	expected := "inner row error"
	err := consume(r, func(s []string) error {
		return nil
	}, func(s []string) error {
		return errors.New(expected)
	})

	if err == nil {
		t.Fatal("consume did not return an error")
	}

	inner := errors.Unwrap(err)
	if inner == nil {
		t.Fatal("incorrect error: could not unwrap")
	}
	if inner.Error() != expected {
		t.Fatalf("incorrect error: expected %v, got %v", expected, inner.Error())
	}
}

func TestConsumeWithConsumer(t *testing.T) {
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
		c := NewConsumer(func(s []string) error {
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
		c.consume(r)

		if calls != len(testCase.ExpectedRows) {
			t.Fatalf("%v: incorrect num calls to row processing function, expected %v, got %v", testCase.Name, len(testCase.ExpectedRows), calls)
		}
	}
}
