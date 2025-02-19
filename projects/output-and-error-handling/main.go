package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	maxWaitSeconds = 5
	defaultRetry   = 1 * time.Second
)

func handleRetryAfter(retryAfter string) (time.Duration, error) {
	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		waitDuration := time.Duration(seconds) * time.Second

		if seconds > maxWaitSeconds {
			return 0, fmt.Errorf("Server requested too long of wait time: %v", seconds)
		}

		return waitDuration, nil
	}

	if datetime, err := time.Parse(http.TimeFormat, retryAfter); err == nil {
		waitDuration := time.Until(datetime)
		if waitDuration <= 0 {
			return 0, fmt.Errorf("retry time is in the past")
		}

		if waitDuration > maxWaitSeconds*time.Second {
			return 0, fmt.Errorf("sever request too long of a wait: %v", waitDuration)
		}

		return waitDuration, nil
	}

	fmt.Fprintf(os.Stderr, "Count not parse Retry-After header. Using deafult retry of %v\n", defaultRetry)

	return defaultRetry, nil

}

func getWeatherData(url string) (string, error) {

	client := &http.Client{}

	res, err := client.Get(url)

	if err != nil {
		return "", fmt.Errorf("Failed to fetch whether data: %w", err)
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read weather response: %w", err)
		}
		return string(body), nil

	case 429:
		retryAfter := res.Header.Get("Retry-After")
		waitDuration, err := handleRetryAfter(retryAfter)
		if err != nil {
			return "", fmt.Errorf("cannot get weather: %w", err)
		}

		if waitDuration > time.Second {
			fmt.Fprintf(os.Stderr, "This might take a while - waiting %v before retrying...\n", waitDuration)
		}

		time.Sleep(waitDuration)
		return getWeatherData(url)

	default:
		return "", fmt.Errorf("unexpected error occured: %d", res.StatusCode)
	}
}

func main() {
	data, err := getWeatherData("http://localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v \n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, data)
}
