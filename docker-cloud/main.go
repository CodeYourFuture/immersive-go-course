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
		port = "8080"
	}

	config := Config{
		Port: port,
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))
}
