package main

import (
	"io"
)

type FilteringPipe struct {
	writer io.Writer
}

func NewFilteringPipe(writer io.Writer) FilteringPipe {
	return FilteringPipe{
		writer: writer,
	}
}

func (fp *FilteringPipe) Write(bytes []byte) (int, error) {
	for i := range bytes {
		if bytes[i] < '0' || bytes[i] > '9' {
			if _, err := fp.writer.Write(bytes[i : i+1]); err != nil {
				return i, err
			}
		}
	}
	// We return len(bytes) because io.Writer is documented to return how many bytes were processed, not how many were actually used.
	return len(bytes), nil
}
