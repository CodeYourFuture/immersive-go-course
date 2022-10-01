package main

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Consumer[Ctx any] struct {
	Header func([]string, Ctx) error
	Row    func([]string, Ctx) error
}

func NewConsumer[Ctx any](headerF func([]string, Ctx) error, rowF func([]string, Ctx) error) *Consumer[Ctx] {
	return &Consumer[Ctx]{
		Header: headerF,
		Row:    rowF,
	}
}

// Consume each row of the CSV, running a function on the header row and each row in turn.
func (c *Consumer[Ctx]) consume(r *csv.Reader, context Ctx) error {
	for {
		// Read a line
		row, err := r.Read()
		// If it's the end of the file, break out of the loop
		if err == io.EOF {
			break
		}
		// Some other error is a problem, so we should return it
		if err != nil {
			// Wrap the error using %w
			return fmt.Errorf("could not read row: %w", err)
		}
		// r.Read keeps track of where we are in the file, so we use that
		if line, _ := r.FieldPos(0); line == 1 {
			// Process the header row
			err = c.Header(row, context)
			if err != nil {
				return fmt.Errorf("header error: %w", err)
			}
			continue
		}
		// Process a non-header row
		err = c.Row(row, context)
		if err != nil {
			return fmt.Errorf("row error: %w", err)
		}
	}
	return nil
}

func consume(r *csv.Reader, headerF func([]string, any) error, rowF func([]string, any) error) error {
	c := NewConsumer(headerF, rowF)
	return c.consume(r, nil)
}

// Consume the channel, returning a slice of the output values
func chanToSlice[T any](in chan T) []T {
	out := make([]T, 0, 100)
	for v := range in {
		out = append(out, v)
	}
	return out
}
