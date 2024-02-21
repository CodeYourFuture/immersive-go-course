package main

import (
	"flag"
	"log"
	"path/filepath"
	"servers/static"
)

func main() {
	path := flag.String("path", ".", "path to static files")
	port := flag.Int("port", 8082, "port the server will listen on")
	flag.Parse()

	absPath, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatal(static.Run(static.Config{
		Dir:  absPath,
		Port: *port,
	}))
}
