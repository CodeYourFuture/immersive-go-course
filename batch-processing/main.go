package main

import (
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

func main() {
	rows := []Row{
		{"files/first"},
		{"files/second"},
		{"files/third"},
	}

	urls := make(chan string, len(rows))
	downloads := make(chan DownloadResult, len(rows))
	processes := make(chan ProcessResult, len(rows))
	uploads := make(chan UploadResult, len(rows))

	for _, row := range rows {
		urls <- row.Url
	}

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(rows))
		go func() {
			wg.Wait()
			close(downloads)
		}()
		for url := range urls {
			log.Printf("loaded: %s\n", url)
			go func(url string) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
				downloads <- DownloadResult{
					Url:      url,
					Filepath: url,
					Error:    nil,
				}
			}(url)
		}
	}()

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(rows))
		go func() {
			wg.Wait()
			close(processes)
		}()
		for dR := range downloads {
			log.Printf("downloaded: %s\n", dR.Filepath)
			go func(dR DownloadResult) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(1000)+100) * time.Millisecond)
				processes <- ProcessResult{
					DownloadResult: dR,
					Filepath:       dR.Filepath,
					Error:          nil,
				}
			}(dR)
		}
	}()

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(rows))
		go func() {
			wg.Wait()
			close(uploads)
		}()
		for pR := range processes {
			log.Printf("processed: %s\n", pR.Filepath)
			go func(pR ProcessResult) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
				uploads <- UploadResult{
					ProcessResult: pR,
					Url:           pR.Filepath,
					Error:         nil,
				}
			}(pR)
		}
	}()

	for uR := range uploads {
		log.Printf("uploaded: %s\n", uR.Url)
	}
}
