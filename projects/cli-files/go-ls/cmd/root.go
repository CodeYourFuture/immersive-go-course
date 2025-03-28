package cmd

import (
	"fmt"
	"log"
	"os"
)
const defaultRoute = "."

func Execute() {
    file := defaultRoute

    if len(os.Args) == 2 {
        file = os.Args[1]

        if file == "-h" {
            fmt.Fprint(os.Stdout, "how to use the go-ls command:\n * type go-ls and a valid directory file\n")
            return
        }

    } else if len(os.Args) > 2 {
        fmt.Fprint(os.Stderr, "this is your current directory. please type just one arg\n")
    }


    FileInfo, err := os.Stat(file)
    if err != nil {
        log.Fatal(err)
    }
    if !FileInfo.IsDir() {
        fmt.Fprintf(os.Stderr, "File %s is not a directory", FileInfo.Name())
        return
    }

    files, err := os.ReadDir(file)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        fmt.Fprintln(os.Stdout, file.Name())
    }
}
