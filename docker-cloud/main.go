package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)


func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world."))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(([]byte("Hello!")))
}

func main() {

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "80"
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ping", pingHandler)

	fmt.Println("Listening")
	log.Fatal(http.ListenAndServe(":" + httpPort, nil))
}

