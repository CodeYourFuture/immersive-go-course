package main

import (
	"encoding/csv"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

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
	// TODO: We don't need to read all here. We can read line-by-line and put each URL to the channel.
	//       This would come at the end so that we don't have to create a buffered channel.
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Skip the header row of the CSV
	// TODO: validate the CSV file
	rows := records[1:]
	// Create an initial input channel for the URLs from each (simulated) CSV row.
	// Buffer the channel so that we can load it up even with no takers.
	urls := make(chan string, len(rows))
	// Push each input into the channel...
	for _, row := range rows {
		urls <- row[0]
	}
	// ...and then immediately close it as we know we have nothing more to add.
	// Takers will be able to take from the channel until the buffer is empty, and then they'll
	// see the closed value.
	close(urls)

	// For each URL, download the file, and pass the path to the next step.
	downloads := Map(urls, download)

	// For each downloaded file, process the image, write a new file and pass it on.
	processes := Map(downloads, process)

	// For each processes file, upload the image, and pass the resulting URL on.
	uploads := Map(processes, upload)

	// Iterate through the completed uploads, logging the final URL.
	for uR := range uploads {
		log.Printf("uploaded: %s\n", uR.Url)
	}
}
