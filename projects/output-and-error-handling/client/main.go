package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		fmt.Println("Sorry we cannot get you the weather")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	switch resp.StatusCode {
	case 200:
		fmt.Fprint(os.Stdout, sb+"\n")
	case 429:
		retrySeconds := 0
		retryTime := resp.Header.Get("Retry-After")

		if retryTime == "a while" {
			fmt.Println("We will retry to get you the weather. Please wait 3 seconds")
			time.Sleep(3 * time.Second)
			main()
		}

		retryTimeDate, err := time.Parse(time.RFC1123, retryTime)

		if err == nil {
			retrySeconds = int(time.Until(retryTimeDate).Seconds())
		} else {
			retrySeconds, err = strconv.Atoi(retryTime)
			if err != nil {
				log.Fatalf("Internal error")
			}
		}
		if retrySeconds > 1 && retrySeconds <= 5 {
			fmt.Printf("We will retry to get you the weather. Please wait %d seconds\n", retrySeconds)
			time.Sleep(time.Duration(retrySeconds) * time.Second)
			main()
		} else {
			fmt.Println("Sorry we cannot get you the weather")
		}
	case 500:
		fmt.Fprint(os.Stderr, sb+"\n")
	default:
		fmt.Println("Sorry we cannot get you the weather")
	}

}
