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
	// Connect to the server and getting a response
	resp, err := http.Get("http://localhost:8080")
	// Show error message if connection is not established
	if err != nil {
		handleError()
	}

	// Close the response body after been fully read
	defer resp.Body.Close()

	// Read response body and store it in a var
	body, err := io.ReadAll(resp.Body)

	// Show error message if we cannot read the response body
	if err != nil {
		handleError()
	}
	// Convert response body from binary to string
	sb := string(body)

	// Handle cases depending on the Status code of the response
	switch resp.StatusCode {
	case 200:
		fmt.Fprint(os.Stdout, sb+"\n")
	case 429:
		handleRateLimited(resp.Header.Get("Retry-After"))
	case 500:
		fmt.Fprint(os.Stderr, sb+"\n")
	default:
		handleError()
	}
}

func handleError() {
	fmt.Println("Sorry we cannot get the weather")
}

// Handle response and retry depending on the Retry-After header
func handleRateLimited(retryTime string) {
	retrySeconds := 0
	var err = error(nil)
	retryTimeDate, err := time.Parse(time.RFC1123, retryTime)

	if retryTime == "a while" {
		retrySeconds = 3
	} else if err == nil {
		retrySeconds = int(time.Until(retryTimeDate).Seconds())
	} else {
		retrySeconds, err = strconv.Atoi(retryTime)
		if err != nil {
			handleError()
		}
	}
	if retrySeconds > 1 && retrySeconds <= 5 {
		fmt.Printf("We will retry to get you the weather. Please wait %d seconds\n", retrySeconds)
		time.Sleep(time.Duration(retrySeconds) * time.Second)
		main()
	} else {
		handleError()
	}
}
