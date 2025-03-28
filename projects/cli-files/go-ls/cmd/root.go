package cmd

import (
	"fmt"
	"log"
	"os"
)
const defaultRoute = "."

func Execute() {
    file := defaultRoute
    
    // Display instructions.
    if len(os.Args) == 2 {
        if os.Args[1] == "-h" {
            fmt.Fprint(os.Stdout, "how to use the go-ls command:\n * type go-ls and a valid directory file\n * type go-ls -m and a valid directory for a better format")
            return
        }
        
        // Read under normal formatting.
        file = os.Args[1]
        if err := readAndPrintFile(file, false); err != nil {
            fmt.Fprintf(os.Stderr, "There was an error reading the directory")
        }
        return
    } 

    // Check valid conditions.
    if len(os.Args) > 3 {
        fmt.Fprint(os.Stderr, "this is your current directory. please type just one arg\n")
        return
    }

    // Read under specical fomatting.
    if len(os.Args) == 3 && os.Args[1] == "-m" {
        file = os.Args[2]
            if checkNonDir(file) {
                if err := readAndPrintFile(file, true); err != nil {
                    fmt.Fprintf(os.Stderr, "There was an error reading the directory")
            }
        }
        return
    }
    if err := readAndPrintFile(file, false); err != nil {
        fmt.Fprintf(os.Stderr, "There was an error reading the directory")
    }
    return
}

func readAndPrintFile(file string, formatted bool) error {
    files, err := os.ReadDir(file)
    if err != nil {
        return err
    }
    for _, file := range files {
        if formatted {
            fmt.Fprintf(os.Stdout, "%s ",  file.Name())
            continue
        }
        fmt.Fprintln(os.Stdout, file.Name())
    }
    return nil
}

// Evaluate possible non-dir.
func checkNonDir(file string) bool {
    FileInfo, err := os.Stat(file)
    if err != nil {
        log.Fatal(err)
    }
    if !FileInfo.IsDir() {
        fmt.Fprintf(os.Stderr, "File %s is not a directory", FileInfo.Name())
        return false
    }
    return true
}
