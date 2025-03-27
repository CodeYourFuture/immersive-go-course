package cmd

import (
	"fmt"
	"log"
	"os"
)
const defaultRoute = "."

func Execute() {
    file := defaultRoute

    if len(os.Args) >= 1 {
       file = os.Args[1]
    }

    files, err := os.ReadDir(file)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        fmt.Println(file)
    }
}
