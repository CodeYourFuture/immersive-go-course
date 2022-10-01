package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

type Context struct {
	Input chan string
}

func consumeHeader(headerRow []string, ctx Context) error {
	if len(headerRow) != 1 {
		return fmt.Errorf("CSV header is incorrect length, expected 1, got %d", len(headerRow))
	}
	if headerRow[0] != "url" {
		return fmt.Errorf("CSV header is incorrect, expected \"url\", got \"%s\"", headerRow[0])
	}
	return nil
}

func consumeRow(row []string, ctx Context) error {
	// TODO: validate row
	log.Printf("url: %v", row[0])
	ctx.Input <- row[0]
	return nil
}

func main() {
	// We need a file to read from
	file := flag.String("file", "", "A path to a CSV with URLs of images to be processed")
	flag.Parse()
	if *file == "" {
		log.Fatal("supply a file using --file")
	}

	// Open the file supplied
	in, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}

	// Read the file using the encoding/csv package
	r := csv.NewReader(in)

	// Create an initial input channel for the URLs from each CSV row.
	ctx := Context{
		Input: make(chan string),
	}

	// Set up the processing pipeline
	uploads := Map(ctx.Input, pipeline)

	// Build a Consumer that will process the header and rows from the CSV.
	// The (implied) Context for the consumer contains the Input channel that the processing pipeline
	// will take values from.
	c := NewConsumer(consumeHeader, consumeRow)

	// Consume the CSV, pushing each input into the channel...
	err = c.consume(r, ctx)
	if err != nil {
		log.Fatalf("could not read CSV: %v", err)
	}

	// ...and then immediately close it as we know we have nothing more to add.
	// Takers will be able to take from the channel until the buffer is empty, and then they'll
	// see the closed value.
	close(ctx.Input)

	// Iterate through the completed uploads, logging the final URLs.
	uploadResults := chanToSlice(uploads)
	log.Printf("output: %v", uploadResults)
}
