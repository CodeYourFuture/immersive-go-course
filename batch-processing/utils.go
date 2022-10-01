package main

import (
	"encoding/csv"
	"fmt"
	"io"
)

// Consume each row of the CSV, running a function on the header row and each row in turn.
func consume(r *csv.Reader, headerF func([]string) error, rowF func([]string) error) error {
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
			err = headerF(row)
			if err != nil {
				return fmt.Errorf("header error: %w", err)
			}
			continue
		}
		// Process a non-header row
		err = rowF(row)
		if err != nil {
			return fmt.Errorf("row error: %w", err)
		}
	}
	return nil
}
