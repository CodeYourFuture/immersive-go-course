package main

import (
	"flag"
	"log"
	"os"
	"servers/api"
)

func main() {
	// Check that DATABASE_URL is set
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatalln("DATABASE_URL not set")
	}

	port := flag.Int("port", 8081, "port the server will listen on")
	flag.Parse()

	log.Fatal(api.Run(api.Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        *port,
	}))
}
