package cmd

import (
	"fmt"
	"os"
    "errors"
)
const defaultRoute = "."

func Execute() {
    file := defaultRoute
    
    // Display instructions.
    if len(os.Args) == 2 {
        if os.Args[1] == "-h" {
            fmt.Fprintln(os.Stdout, "how to use the go-ls command:\n * type go-ls and a valid directory file\n * type go-ls -m and a valid directory for a better format")
            return
        }
        
        // Read under normal formatting.
        file = os.Args[1]
        if err := readAndPrintFile(file, false); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
        return
    } 

    // Check valid conditions.
    if len(os.Args) > 3 {
        fmt.Fprintln(os.Stderr, "type go-ls -h to chech the commands")
        return
    }

    // Read under specical fomatting.
    if len(os.Args) == 3 && os.Args[1] == "-m" {
        file = os.Args[2]
        if err := readAndPrintFile(file, true); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
        return
    }
    
    // Read default route.
    if err := readAndPrintFile(file, false); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    return
}

func readAndPrintFile(file string, formatted bool) error {
    if err := checkNonDir(file); err != nil {
        return err
    }

    files, err := os.ReadDir(file)
    if err != nil {
        return errors.New("there was an error reading the directory")
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
func checkNonDir(file string) error {
    FileInfo, err := os.Stat(file)
    if err != nil {
        return errors.New("there was an error reading the directory")
    }
    if !FileInfo.IsDir() {
        return errors.New(FileInfo.Name() + " is a file not a directory")
    }
    return nil
}
