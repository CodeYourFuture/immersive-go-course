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
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	resp, err := c.Get("http://localhost:8080")

	if err != nil {
		fmt.Fprint(os.Stderr, "Something went wrong. Cannot complete request at the moment\n")
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 429:
		retryTime := resp.Header.Get("Retry-After")
		formattedRetryTime, err := parseRetryValue(retryTime)

		if err != nil {
			// cannot determine how long to sleep for, retry once after 1 sec delay. Reason - The server response shows it is safe to retry. However we might not want keep at it for long - hence the choice

			sleep(1)
			resp, err = retry(1, c)

			if err != nil {
				// retry still unsuccessful and server does not respond
				fmt.Fprint(os.Stderr, "Something went wrong. Cannot complete request at the moment\n")
				resp.Body.Close()
				return
			}

		} else if formattedRetryTime > 5 {
			// sleep time is long
			fmt.Fprint(os.Stderr, "Server busy. Cannot retrieve weather\n")
			resp.Body.Close()
			return
		} else {
			sleep(formattedRetryTime)
			resp, err = retry(3, c) // chose 3 as max number of retries

			if err != nil {
				// retry still unsuccessful and server does not respond
				fmt.Fprint(os.Stderr, "Something went wrong. Cannot complete request at the moment\n")
				resp.Body.Close()
				return
			}

		}
	}

	// for other responses - print to stdout
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(os.Stdout, string(body), "\n")
}

func parseRetryValue(v string) (int, error) {

	value, err := strconv.Atoi(v)
	if err != nil {
		// check if retry time is in httpTime format
		httpTime, httpTimeErr := time.Parse(time.RFC1123, v)

		if httpTimeErr == nil {
			value = int(time.Until(httpTime).Seconds())
			err = httpTimeErr
		}
	}
	return value, err

}

func retry(maxRetries int, c http.Client) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 1; i <= maxRetries; i++ {
		// max of three retries after which we will give up if not successful
		fmt.Fprint(os.Stderr, "Retrying ..........\n")

		time.Sleep(time.Duration(1) * time.Second) // delay between each retries

		resp, err = c.Get("http://localhost:8080")
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 {
			break
		}

	}
	return resp, err
}

// sleep for a given time
func sleep(t int) {
	if t > 1 {
		fmt.Fprint(os.Stderr, "We are currently retrying your request. Things might take a bit longer than usual\n")
	}
	time.Sleep(time.Duration(t) * time.Second)
}
