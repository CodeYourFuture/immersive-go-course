package main

import (
	"net/http"
	"os"

	"github.com/CodeYourFuture/immersive-go-course/projects/output-and-error-handling/client"
)

func main() {
    // Inside main I create a real API
    c := &http.Client{}
    api := client.BaseAPI{
        Client: c,
        URL: "http://127.0.0.1:8080",
        Testing: false,
    }
    api.DoStuff(os.Stdout, os.Stderr)
}
