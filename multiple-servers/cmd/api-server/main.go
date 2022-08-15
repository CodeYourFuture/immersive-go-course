package main

import (
	"log"
	"os"
	"servers/api"
)

func main() {
	// Check that DATABASE_URL is set
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatalln("DATABASE_URL not set")
	}

	log.Fatal(api.Run(api.Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        8081,
	}))
}
