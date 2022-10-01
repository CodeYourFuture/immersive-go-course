package main

import (
	"encoding/csv"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

type Context struct {
	Input chan string
}

type DownloadResult struct {
	Url      string
	Filepath string
	Error    error
}

type ProcessResult struct {
	DownloadResult DownloadResult
	Filepath       string
	Error          error
}

type UploadResult struct {
	ProcessResult ProcessResult
	Url           string
	Error         error
}

func download(url string) DownloadResult {
	log.Printf("downloading: %v\n", url)
	time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
	return DownloadResult{
		Url:      url,
		Filepath: url,
		Error:    nil,
	}
}

func process(dR DownloadResult) ProcessResult {
	log.Printf("processing: %v\n", dR.Filepath)
	time.Sleep(time.Duration(rand.Intn(1000)+100) * time.Millisecond)
	return ProcessResult{
		DownloadResult: dR,
		Filepath:       dR.Filepath,
		Error:          nil,
	}
}

func upload(pR ProcessResult) UploadResult {
	log.Printf("uploading: %v\n", pR.Filepath)
	time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
	return UploadResult{
		ProcessResult: pR,
		Url:           pR.Filepath,
		Error:         nil,
	}
}

func pipeline(ctx Context) chan UploadResult {
	// For each URL, download the file, and pass the path to the next step.
	downloads := Map(ctx.Input, download)

	// For each downloaded file, process the image, write a new file and pass it on.
	processes := Map(downloads, process)

	// For each processes file, upload the image, and pass the resulting URL on.
	return Map(processes, upload)
}

func consumeHeader(headerRow []string, ctx Context) error {
	// TODO: validate header
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

	// Create an initial input channel for the URLs from each (simulated) CSV row.
	ctx := Context{
		Input: make(chan string),
	}

	// Set up the processing pipeline
	uploads := pipeline(ctx)

	// Build a Consumer that will process the header and rows from the CSV.
	// The (implied) Context for the consumer contains the Input channel that the processing pipeline
	// will take values from.
	c := NewConsumer(consumeHeader, consumeRow)

	// Consume the CSV, pushing each input into the channel...
	c.consume(r, ctx)

	// ...and then immediately close it as we know we have nothing more to add.
	// Takers will be able to take from the channel until the buffer is empty, and then they'll
	// see the closed value.
	close(ctx.Input)

	// Iterate through the completed uploads, logging the final URLs.
	uploadResults := chanToSlice(uploads)
	log.Printf("output: %v", uploadResults)
}
