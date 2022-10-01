package main

import (
	"log"
	"math/rand"
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
	if dR.Error != nil {
		return ProcessResult{
			DownloadResult: dR,
			Filepath:       "",
			Error:          dR.Error,
		}
	}
	log.Printf("processing: %v\n", dR.Filepath)
	time.Sleep(time.Duration(rand.Intn(1000)+100) * time.Millisecond)
	return ProcessResult{
		DownloadResult: dR,
		Filepath:       dR.Filepath,
		Error:          nil,
	}
}

func upload(pR ProcessResult) UploadResult {
	if pR.Error != nil {
		return UploadResult{
			ProcessResult: pR,
			Url:           "",
			Error:         pR.Error,
		}
	}
	log.Printf("uploading: %v\n", pR.Filepath)
	time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
	return UploadResult{
		ProcessResult: pR,
		Url:           pR.Filepath,
		Error:         nil,
	}
}

func pipeline(url string) UploadResult {
	dR := download(url)
	pR := process(dR)
	uR := upload(pR)
	return uR
}
