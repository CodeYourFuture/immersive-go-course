package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	client := &http.Client{Timeout: time.Duration(1) * time.Second}

	response, err := client.Get("http://localhost:8080")
	if err != nil {
		fmt.Fprint(os.Stderr, "Server is down. Please try again later\t")
		os.Exit(1)
	}

	defer response.Body.Close()
	
	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Fprint(os.Stderr, "Response body could not be read: ", err)
			os.Exit(2)
		}
		fmt.Fprintln(os.Stdout, string(body))
	} else if response.StatusCode == http.StatusTooManyRequests {
		handleRetry(response)
	} else {
		fmt.Fprintf(os.Stderr, "Unexpected response: %d\n", response.StatusCode)
		os.Exit(4)
	}
	
	
}

func handleRetry(res *http.Response) {

	retryHeader := res.Header.Get("Retry-After")
	parsedTime, err := http.ParseTime(retryHeader)
	if err == nil {
		waitTime := time.Until(parsedTime)
		fmt.Printf("You have to wait for %vsecs to restart the application", int64(waitTime/time.Second))
		time.Sleep(waitTime)
	} else {
		waitSecs, err := strconv.Atoi(retryHeader)
		if err == nil {
			fmt.Printf("You have to wait for %dsecs to start the application again", waitSecs)
			time.Sleep(time.Duration(waitSecs) * time.Second)
		} else {
			fmt.Fprint(os.Stderr, "Invalid Retry Header!!!\t")
			os.Exit(3)
		}
	}

	defer res.Body.Close()
}
