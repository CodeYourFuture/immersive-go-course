package main

import (
	"fmt"
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
