package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port string
}

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "80"
	}

	config := Config{
		Port: port,
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world."))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))
}
