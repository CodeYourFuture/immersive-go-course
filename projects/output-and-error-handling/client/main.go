package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func fetch(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main() {
	resp, err := fetch("http://localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	switch resp.StatusCode {
	case 200:
		fmt.Fprint(os.Stdout, sb)
	case 429:
		fmt.Print(sb)
	case 500:
		fmt.Print(sb)
	}
}
