package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Port int
}

func main() {
	config := Config{
		Port: 8090,
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
