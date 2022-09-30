package main

import (
	"log"
	"math/rand"
	"time"
)

type Row struct {
	Url string
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

func main() {
	rows := []Row{
		{"files/first"},
		{"files/second"},
		{"files/third"},
	}

	// Create an initial input channel for the URLs from each (simulated) CSV row.
	// Buffer the channel so that we can load it up even with no takers.
	urls := make(chan string, len(rows))
	// Push each input into the channel...
	for _, row := range rows {
		urls <- row.Url
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
