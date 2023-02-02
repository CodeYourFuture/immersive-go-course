package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// We generate a random number between 0 and 9 (inclusive), so that we can decide whether to behave properly (half of the time), or simulate error conditions.
		randomNumber := rand.Intn(10)
		if randomNumber < 5 {
			// 50% of the time, we just report the weather.
			if randomNumber < 3 {
				w.Write([]byte("Today it will be sunny!\n"))
			} else {
				w.Write([]byte("I'd bring an umbrella, just in case...\n"))
			}
		} else if randomNumber < 8 {
			// 30% of the time, we say we're too busy and say try again in a few seconds.

			// Generate a random number between 1 and 10, for the number of seconds to tell the client to wait before retrying:
			retryAfterSeconds := rand.Intn(9) + 1

			// 10% of the time we give a number of seconds to wait.
			retryAfter := strconv.Itoa(retryAfterSeconds)
			if randomNumber == 6 {
				// 10% of the time we give a timestamp to wait until.
				timeAfterDelay := time.Now().UTC().Add(time.Duration(retryAfterSeconds) * time.Second)
				retryAfter = timeAfterDelay.Format(http.TimeFormat)
			} else if randomNumber == 7 {
				// But 10% of the time there's actually a bug which means we don't tell you a time to retry after, and trying to parse the header will result in an error.
				retryAfter = "a while"
			}
			w.Header().Set("Retry-After", retryAfter)
			w.WriteHeader(429)
			w.Write([]byte("Sorry, I'm too busy"))
		} else {
			// 20% of the time we just drop the connection.
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
		}
	})

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}
}
