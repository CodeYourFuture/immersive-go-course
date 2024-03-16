package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const serverPort = 8080

func main() {
	requestURL := fmt.Sprintf("http://localhost:%d", serverPort)
	resp, err := http.Get(requestURL)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)

	if resp.StatusCode == 200 {
		fmt.Fprint(os.Stdout, sb)
	}

}
