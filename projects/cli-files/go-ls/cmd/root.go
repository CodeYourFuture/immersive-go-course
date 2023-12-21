package cmd

import (
	"fmt"
	"os"
)

func Execute(dir string) {
	fileInfo, _ := os.Stat(dir)
	if !fileInfo.IsDir() {
		fmt.Println(fileInfo.Name())
		return
	}

	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
