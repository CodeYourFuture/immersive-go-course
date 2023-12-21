package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type retryableError struct {
	retryAfter time.Duration
}

func (err retryableError) Error() string {
	return fmt.Sprintf("retryable error: retry after %v", err.retryAfter)
}

func main() {
	const endpoint = "http://localhost:8080"
	const maxRetries = 5

	for attempt := 0; attempt < maxRetries; attempt++ {
		weather, err := requestWeather(endpoint)
		if err == nil {
			fmt.Println(weather)
			return
		}

		var retryableErr retryableError
		if errors.As(err, &retryableErr) {
			if retryableErr.retryAfter > 5*time.Second {
				log.Fatal("Can't give you the weather info. Retry-After exceeds limit.")
			} else if retryableErr.retryAfter > time.Second {
				log.Println("Giving you the weather info shortly.")
			}
			time.Sleep(retryableErr.retryAfter)
		} else {
			log.Fatal("Can't give you the weather info. Unexpected communication error.")
		}
	}

	log.Fatal("Exceeded maximum retry attempts. Can't give you the weather info.")
}

func requestWeather(endpoint string) (string, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("http get error requesting the weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		retryAfterHeader := resp.Header.Get("Retry-After")
		seconds, err := strconv.Atoi(retryAfterHeader)
		if err == nil {
			return "", retryableError{retryAfter: time.Duration(seconds) * time.Second}
		}

		parsedTime, err := http.ParseTime(retryAfterHeader)
		if err == nil {
			return "", retryableError{retryAfter: time.Until(parsedTime)}
		}

		return "", fmt.Errorf("error parsing the retry after header")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing the response body: %w", err)
	}
	return string(body), nil
}
