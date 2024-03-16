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
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		handleError()
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		handleError()
	}

	sb := string(body)

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
	fmt.Println("Sorry we cannot get you the weather")
}

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
