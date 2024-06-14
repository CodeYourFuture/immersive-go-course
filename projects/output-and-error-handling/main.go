package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/CodeYourFuture/immersive-go-course/projects/output-and-error-handling/fetcher"
)

func main() {
	f := fetcher.WeatherFetcher{}
	// Loop because we may need to retry.
	for {
		if weather, err := f.Fetch("http://localhost:8080/"); err != nil {
			// If we're told to retry, do so.
			if errors.Is(err, fetcher.ErrRetry) {
				continue
			}
			// Otherwise tell the user there was an error and give up.
			fmt.Fprintf(os.Stderr, "Error getting weather: %v\n", err)
			os.Exit(1)
		} else {
			// Print out the weather and be happy.
			fmt.Fprintf(os.Stdout, "%s\n", weather)
			break
		}
	}
}
