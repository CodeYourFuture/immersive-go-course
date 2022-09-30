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

func main() {
	rows := []Row{
		{"files/first"},
		{"files/second"},
		{"files/third"},
	}

	urls := make(chan string, len(rows))
	for _, row := range rows {
		urls <- row.Url
	}
	close(urls)

	downloads := Map(
		urls,
		func(url string) DownloadResult {
			log.Printf("downloading: %v\n", url)
			time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
			return DownloadResult{
				Url:      url,
				Filepath: url,
				Error:    nil,
			}
		},
	)

	processes := Map(
		downloads,
		func(dR DownloadResult) ProcessResult {
			log.Printf("processing: %v\n", dR.Filepath)
			time.Sleep(time.Duration(rand.Intn(1000)+100) * time.Millisecond)
			return ProcessResult{
				DownloadResult: dR,
				Filepath:       dR.Filepath,
				Error:          nil,
			}
		},
	)

	uploads := Map(
		processes,
		func(pR ProcessResult) UploadResult {
			log.Printf("uploading: %v\n", pR.Filepath)
			time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
			return UploadResult{
				ProcessResult: pR,
				Url:           pR.Filepath,
				Error:         nil,
			}
		},
	)

	for uR := range uploads {
		log.Printf("uploaded: %s\n", uR.Url)
	}
}
