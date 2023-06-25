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

var httpClient = http.Client{
	Timeout: time.Second * 30,
}

func main() {
	for {
		resp, err := get()
		if err != nil {
			writeStdErr(err.Error())
			return
		}

		switch {
		case resp.StatusCode == 200:
			body, _ := io.ReadAll(resp.Body)
			writeStdOut(string(body) + "\n")
			return
		case resp.StatusCode == 429:
			duration, err := parseRetryHeader(resp.Header.Get(retryAfter))
			if err != nil {
				break
			}
			writeStdOut(fmt.Sprintf("waiting for: %v\n", duration))
			time.Sleep(duration)
			continue
		default:
			break
		}
	}
}

func parseRetryHeader(retry string) (time.Duration, error) {
	wait, err := time.Parse(http.TimeFormat, retry)
	switch {
	case err == nil:
		return wait.Sub(time.Now()), nil
	default:
		wait, err := strconv.Atoi(retry)
		if err != nil {
			return 0, err
		}
		return time.Duration(wait), nil
	}
}

func get() (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, serverURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating a new request: %w\n", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection to the server failed: %w\n", err)
	}

	return resp, nil
}

func writeStdOut(format string) {
	write(os.Stdout, format)
}

func writeStdErr(format string) {
	write(os.Stderr, format)
}

func write(w io.Writer, format string) {
	if _, err := fmt.Fprintf(w, format); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("error writing to console: %v\n", err))
		return
	}
}
