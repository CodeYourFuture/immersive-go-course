package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Global constants
const (
	delayLimit = 5
	url        = "http://localhost:8080"
)

func main() {

	// Fetch weather updates
	update, err := getWeatherUpdate(url)

	// Handle errors and
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v ", err)
		os.Exit(1)
	} else if update == "" {
		fmt.Fprint(os.Stderr, "error: no updates for today ")
		os.Exit(1)
	}

	// Prints weather update
	fmt.Fprintln(os.Stdout, update)

}

// Function to fetch weather updates
func getWeatherUpdate(url string) (string, error) {

	// HTTP GET request
	resp, err := http.Get(url)

	if err != nil {
		return "", errors.New("sorry, could not get weather updates from " + url)
	}

	defer resp.Body.Close()

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK:
		return parseBody(resp)
	case http.StatusTooManyRequests:
		return "", retryServer(resp)
	default:
		return "", errors.New("an unexpected error occured, try later ")
	}

}

// Function to parse response body
func parseBody(request *http.Response) (string, error) {

	body, err := io.ReadAll(request.Body)

	if err != nil {
		return "", errors.New("could not get the weather response ")
	}

	weatherUpdate := string(body)

	return weatherUpdate, nil

}

// Function to handle retry logic based on Retry-After header
func retryServer(request *http.Response) error {

	const retryHeader = "Retry-After"
	retryTime := request.Header.Get(retryHeader)

	delay, err := strconv.Atoi(retryTime)

	if err != nil {
		return errors.New("service temporarily unavailable, retry later ")
	}

	// if delay within limit, retry after waiting
	if delay <= delayLimit {
		fmt.Fprintln(os.Stderr, "retrying, it might a take a while.....")
		time.Sleep(time.Duration(delay) * time.Second)
		getWeatherUpdate(url) //Retry fetching weather updates
		return nil
	}

	return nil
}
