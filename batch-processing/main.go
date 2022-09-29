package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
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

func pipe[In any, Out any](in <-chan In, work func(in In) Out) chan Out {
	out := make(chan Out)
	go func() {
		for inV := range in {
			go func(inV In) {
				out <- work(inV)
			}(inV)
		}
		fmt.Printf("closing %v", out)
		close(out)
	}()
	return out
}

func main() {
	rows := []Row{
		{"files/first"},
		{"files/second"},
		{"files/third"},
	}

	urls := make(chan string, len(rows))
	defer close(urls)

	for _, row := range rows {
		urls <- row.Url
	}

	downloads := pipe(
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

	processes := pipe(
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

	uploads := pipe(
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

	var wg sync.WaitGroup
	wg.Add(len(rows))
	go func() {
		for uR := range uploads {
			wg.Done()
			log.Printf("uploaded: %s\n", uR.Url)
		}
	}()
	wg.Wait()
}
