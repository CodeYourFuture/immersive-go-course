package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
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

func genFilepath(suffix string) string {
	return fmt.Sprintf("/tmp/%d-%d.%s", time.Now().UnixMilli(), rand.Int(), suffix)
}

func download(url string) DownloadResult {
	log.Printf("downloading: %v\n", url)
	filepath := genFilepath("jpg")

	// Create a new file that we will write to
	out, err := os.Create(filepath)
	if err != nil {
		return DownloadResult{
			Url:      url,
			Filepath: "",
			Error:    err,
		}
	}
	defer out.Close()

	// Get it from the internet!
	res, err := http.Get(url)
	if err != nil {
		return DownloadResult{
			Url:      url,
			Filepath: "",
			Error:    err,
		}
	}
	defer res.Body.Close()

	// TODO: check status code

	// Copy the body of the response to the created file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return DownloadResult{
			Url:      url,
			Filepath: "",
			Error:    err,
		}
	}

	return DownloadResult{
		Url:      url,
		Filepath: filepath,
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
