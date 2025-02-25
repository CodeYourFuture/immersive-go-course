package cmd

import (
	"fmt"
	"os"
	"strings"
)

func listDir(dir string, commaSeparated bool) {
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !info.IsDir() {
		fmt.Println(dir)
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	var output []string
	for _, file := range files {
		output = append(output, file.Name())
	}

	if commaSeparated {
		fmt.Println(strings.Join(output, ", "))
	} else {
		fmt.Println(strings.Join(output, " "))
	}
}

func Execute() {
	args := os.Args[1:]
	commaSeparated := false
	var directories []string

	// Parse arguments correctly
	for _, arg := range args {
		if arg == "-h" {
			fmt.Println(`Usage: go-ls [OPTIONS] [DIRECTORY]...
A simple implementation of "ls" in Go that lists files in a given directory or prints the file name if it is not a directory.

Available options:
  -h  Displays this help message
  -m  List files separated by commas instead of spaces`)
			return
		} else if arg == "-m" {
			commaSeparated = true
		} else {
			// Assume any non-option argument is a directory
			directories = append(directories, arg)
		}
	}

	// If no directories are specified, list the current directory
	if len(directories) == 0 {
		listDir(".", commaSeparated)
		return
	}

	// If more than one directory, print the directory name before listing
	printDirName := len(directories) > 1

	for _, dir := range directories {
		if printDirName {
			fmt.Printf("%s:\n", dir)
		}
		listDir(dir, commaSeparated)
		fmt.Println()
	}
}
