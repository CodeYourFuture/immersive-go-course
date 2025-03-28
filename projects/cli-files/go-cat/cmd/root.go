package cmd

import (
	"fmt"
	"log"
	"os"
)

func Execute() {
    if len(os.Args) == 1 {
        fmt.Fprintln(os.Stderr, "Please type a valid file to cat")
        return
    }

    for _, file := range os.Args[1:] {
        if !isValidFile(file) {
            fmt.Fprintln(os.Stderr, "Please type a valid file to cat")
            return
        }
        buf, err := os.ReadFile(file)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Fprintln(os.Stdout, string(buf))
    }
}

func isValidFile(file string) bool {
    FileInfo, _ := os.Stat(file)

    if FileInfo.IsDir() {
        return false
    }
    return true
}
