package main

import (
	"flag"
	"log"
	"path/filepath"
	"servers/static"
)

func main() {
	path := flag.String("path", ".", "path to static files")
	flag.Parse()

	absPath, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalln(err)
	}

	static.Run(static.Config{
		Dir:  absPath,
		Port: 8080,
	})
}
