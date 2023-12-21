package cmd

import (
	"bufio"
	"fmt"
	"os"
)

type Options struct {
	ShowLineNumbers bool
}

func Execute(input string, opts Options) {
	file, _ := os.Open(input)
	defer file.Close()

	fileInfo, _ := file.Stat()
	if fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "%s: %s: Is a directory\n", os.Args[0], fileInfo.Name())
		return
	}

	scanner := bufio.NewScanner(file)
	if opts.ShowLineNumbers {
		printWithLineNumbers(scanner)
	} else {
		printPlain(scanner)
	}
}

func printPlain(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func printWithLineNumbers(scanner *bufio.Scanner) {
	for i := 1; scanner.Scan(); i++ {
		fmt.Printf("%6d  %s\n", i, scanner.Text())
	}
}
