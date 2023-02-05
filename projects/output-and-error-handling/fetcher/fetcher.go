package fetcher

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ErrRetry = errors.New("should retry")

type WeatherFetcher struct {
	client http.Client
}

// Fetch gets the weather. If it encounters an error, it will either return retryError if the error is retriable, or another error if it is fatal.
func (w *WeatherFetcher) Fetch(url string) (string, error) {
	response, err := w.client.Get(url)
	if err != nil {
		// Add context to the error about what we were trying to do when we encountered it.
		// We don't wrap with something like "couldn't get weather", because our caller is expected to add that kind of context.
		return "", fmt.Errorf("couldn't make HTTP request: %w", err)
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("error trying to read response: %w", err)
		}
		return string(body), nil
	case http.StatusTooManyRequests:
		if err := handle429(response.Header.Get("retry-after")); err != nil {
			return "", fmt.Errorf("error handling 'too many requests' response: %w", err)
		}
		return "", ErrRetry
	default:
		errorDescription := convertHTTPErrorResponseToDescription(response)
		return "", fmt.Errorf("unexpected response from server: %s", errorDescription)
	}
}

func handle429(retryAfterHeader string) error {
	delay, err := parseDelay(retryAfterHeader)
	if err != nil {
		// handle429 is a really small function that doesn't really do much - its job is "parse a header to seconds, then sleep for that many seconds".
		// Accordingly, we don't really have much context to add to the error here, so we won't wrap it.
		return err
	}
	// This code considers each request independently - it would also be very reasonable to keep a timer since when we made the first request, and give up if the _total_ time is going to be more than 5 seconds, rather than the per-request time.
	if delay > 5*time.Second {
		return fmt.Errorf("giving up request: server told us it's going to be too busy for requests for more than the next 5 seconds")
	}
	if delay > 1*time.Second {
		fmt.Fprintf(os.Stderr, "Server reported it's receiving too many requests - waiting %s before retrying\n", delay)
	}
	time.Sleep(delay)
	return nil
}

func parseDelay(retryAfterHeader string) (time.Duration, error) {
	// Try to decode the retry-after header as a whole number of seconds.
	if waitFor, err := strconv.Atoi(retryAfterHeader); err == nil {
		return time.Duration(waitFor) / time.Nanosecond * time.Second, nil
	}
	// If it wasn't a whole number of seconds, maybe it was a timestamp - try to decode that.
	if waitUntil, err := http.ParseTime(retryAfterHeader); err == nil {
		return time.Until(waitUntil), nil
	}
	// If we couldn't parse either of the expected forms of the header, give up.
	// Include the raw value in the error to help with debugging.
	// Note that if this were a web service, we'd probably log the bad value on the server-side, but not return as much information to the user.
	return -1, fmt.Errorf("couldn't parse retry-after header as an integer number of seconds or a date. Value was: %q", retryAfterHeader)
}

func convertHTTPErrorResponseToDescription(response *http.Response) string {
	var bodyString string
	body, err := io.ReadAll(response.Body)
	if err == nil {
		bodyString = string(body)
	} else {
		bodyString = "<error reading body>"
	}
	return fmt.Sprintf("Status code: %s, Body: %s", response.Status, bodyString)
}
