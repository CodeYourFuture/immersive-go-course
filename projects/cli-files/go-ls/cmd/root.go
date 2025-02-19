package cmd

import (
	"fmt"
	"log"
	"os"
)

func Execute() {
	files, err := os.ReadDir(".")
	if err != nil {
		// for loggin we can use log.fatal - which prints to sterr
		// and also exit with 1
		// we can also use printlien which just prints to the stout
		// fmt.Fprintln(os.Stderr, err)
		// os.Exit(1)
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
