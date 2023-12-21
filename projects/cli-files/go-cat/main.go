package main

import (
	"flag"
	"go-cat/cmd"
)

func main() {
	var nFlag bool
	flag.BoolVar(&nFlag, "n", false, "number all output lines")
	flag.Parse()

	opts := cmd.Options{ShowLineNumbers: nFlag}
	for _, input := range flag.Args() {
		cmd.Execute(input, opts)
	}
}
