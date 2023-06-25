package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Headers = string

const (
	serverURL          = "http://localhost:8080"
	retryAfter Headers = "Retry-After"
)

type fetcher struct {
	client http.Client
}

func NewFetcher() *fetcher {
	return &fetcher{
		http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (f *fetcher) fetch() (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, serverURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating a new request: %w\n", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection to the server failed: %w\n", err)
	}

	return resp, nil
}

func main() {
	f := NewFetcher()

	for {
		resp, err := f.fetch()
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("error fetching the weather data: %v\n", err))
			os.Exit(1)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stdout, string(body)+"\n")
			os.Exit(0)
		case http.StatusTooManyRequests:
			duration, err := parseRetryHeader(resp.Header.Get(retryAfter))
			if err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("unexpecte retry error: %v\n", err))
				fmt.Fprintf(os.Stdout, "retrying the request")
				continue
			}
			fmt.Fprintf(os.Stdout, fmt.Sprintf("server asked to wait and retry after: %v\n", duration))
			time.Sleep(duration)
			continue
		default:
			fmt.Fprintf(os.Stderr, "unexpected status code received: %v\n", err)
			os.Exit(1)
		}
	}
}

func parseRetryHeader(retry string) (time.Duration, error) {
	var err error
	wait, err := time.Parse(http.TimeFormat, retry)
	if err == nil {
		return wait.Sub(time.Now()), nil
	}

	waitInSeconds, err := strconv.Atoi(retry)
	if err == nil {
		return time.Duration(waitInSeconds), err
	}

	return -1, fmt.Errorf("unable to parse the retry header: %w", err)
}
