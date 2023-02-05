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
		switch randomNumber {
		// 50% of the time, we just report the weather. 30% nice, 20% less so.
		case 0, 1, 2:
			w.Write([]byte("Today it will be sunny!"))
		case 3, 4:
			w.Write([]byte("I'd bring an umbrella, just in case..."))
		// 30% of the time, we say we're too busy and say try again in a few seconds, in a few different ways.
		case 5:
			// Generate a random number between 1 and 10, for the number of seconds to tell the client to wait before retrying:
			retryAfterSeconds := rand.Intn(9) + 1

			// 10% of the time we give a number of seconds to wait.
			retryAfterHeader := strconv.Itoa(retryAfterSeconds)

			rejectAsTooBusy(w, retryAfterHeader)
		case 6:
			// Generate a random number between 1 and 10, for the number of seconds to tell the client to wait before retrying:
			retryAfterSeconds := rand.Intn(9) + 1

			// 10% of the time we give a timestamp to wait until.
			timeAfterDelay := time.Now().UTC().Add(time.Duration(retryAfterSeconds) * time.Second)
			retryAfterHeader := timeAfterDelay.Format(http.TimeFormat)

			rejectAsTooBusy(w, retryAfterHeader)
		case 7:
			// But 10% of the time there's actually a bug which means we don't tell you a time to retry after, and trying to parse the header will result in an error.
			retryAfter := "a while"
			rejectAsTooBusy(w, retryAfter)
		// 20% of the time we just drop the connection.
		case 8, 9:
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
		default:
			// This shouldn't be possible, as we generated the random number to be at most 9.
			// Print out enough information in our local log that we can notice and debug what happened:
			fmt.Fprintf(os.Stderr, "Reached unreachable code - HTTP handler switch encountered unhandled random number %d which shouldn't be possible", randomNumber)
			// Give a very generic error message to the caller, because they don't know anything about the internals of our code, and we don't want to tell them anything about it.
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
		}
	})

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}
}

func rejectAsTooBusy(w http.ResponseWriter, retryAfterHeader string) {
	w.Header().Set("Retry-After", retryAfterHeader)
	w.WriteHeader(429)
	w.Write([]byte("Sorry, I'm too busy"))
}
