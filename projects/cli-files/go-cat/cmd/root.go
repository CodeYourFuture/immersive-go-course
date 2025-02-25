package cmd

import (
	"fmt"
	"os"
	"strings"
)

func Execute() {
	args := os.Args[1:]
	printLine := false

	for _, arg := range args {

		if arg == "-n" {
			printLine = true
			continue
		}

		info, err := os.Stat(arg)

		if err != nil {
			fmt.Println(err)
			return
		}

		if info.IsDir() {
			fmt.Printf("go-cat: %s: Is a directory \n", arg)
			continue
		}

		data, err := os.ReadFile(arg)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !printLine {
			os.Stdout.Write(data)
			return
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			fmt.Printf("%d %s\n", i, line)
		}

	}
}
