package main

import (
	"go-ls/cmd"
	"log"
	"os"
)

func main() {
	dir, err := parseDir()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Execute(dir)
}

func parseDir() (string, error) {
	if len(os.Args) > 1 {
		return os.Args[1], nil
	}
	return os.Getwd()
}
