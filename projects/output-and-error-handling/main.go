package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	for {
		weather, err := getWeatherWithRetries(3)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(weather)

	
		time.Sleep(1 * time.Second)
	}
}

func getWeatherWithRetries(maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		weather, err := getWeather()
		if err == nil {
			return weather, nil
		}

		fmt.Fprintf(os.Stderr, "Failed to fetch weather: %v\n", err)
		fmt.Fprintf(os.Stderr, "Retrying (%d/%d)...\n", i+1, maxRetries)
		time.Sleep(1 * time.Second) 
	}

	return "", fmt.Errorf("exceeded maximum retries")
}

func getWeather() (string, error) {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		weather := "Weather: " + readBody(resp)
		return weather, nil
	case http.StatusTooManyRequests:
		retryAfterHeader := resp.Header.Get("Retry-After")
		if retryAfterHeader == "" {
			return "", fmt.Errorf("server didn't provide a Retry-After header")
		}

		waitTime, err := parseRetryAfter(retryAfterHeader)
		if err != nil {
			return "", fmt.Errorf("failed to parse Retry-After header: %v", err)
		}

		fmt.Fprintf(os.Stderr, "Server busy, waiting %s before retrying...\n", waitTime)
		time.Sleep(waitTime)

		return getWeather() 
	case http.StatusServiceUnavailable:
		return "", fmt.Errorf("server is temporarily unavailable, giving up")
	default:
		return "", fmt.Errorf("unexpected response from server: %s", resp.Status)
	}
}

func readBody(resp *http.Response) string {
	body := make([]byte, 512) 
	n, _ := resp.Body.Read(body)
	return string(body[:n])
}

func parseRetryAfter(retryAfterHeader string) (time.Duration, error) {
	retryAfterHeader = strings.TrimSpace(retryAfterHeader)

	
	seconds, err := strconv.Atoi(retryAfterHeader)
	if err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	
	t, err := time.Parse(http.TimeFormat, retryAfterHeader)
	if err == nil {
		return t.Sub(time.Now().UTC()), nil
	}

	
	switch strings.ToLower(retryAfterHeader) {
	case "a while":
		return 5 * time.Second, nil
	default:
		return 0, fmt.Errorf("unknown Retry-After value: %s", retryAfterHeader)
	}
}
