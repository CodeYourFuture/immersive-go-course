package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := c.Get("http://localhost:8080")
	if err != nil {
		fmt.Println("Error from request:", err)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response:", string(body))
}
