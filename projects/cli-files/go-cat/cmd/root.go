package cmd

import (
	"bufio"
	"fmt"
	"os"
    "errors"
)

func Execute() {
    if len(os.Args) == 1 {
        fmt.Fprintln(os.Stderr, "Please type a valid file to cat")
        return
    }
   
    // show numbers.
    if os.Args[1] == "-n" {
        iterateOverArgs(2, true, false)
        return
    }

    // show numbers only on non-empty lines.
    if os.Args[1] == "-b" {
        iterateOverArgs(2, true, true)
        return
    }
    
    // default.
    iterateOverArgs(1, false, false) 
    return
}    

func iterateOverArgs(startingIndex int, showNumbers bool, showEspecialNumbers bool) {
    for _, f := range os.Args[startingIndex:] {
        if err := isValidFile(f); err != nil {
            fmt.Fprintln(os.Stderr, err)
            continue
        }
        file, _ := os.Open(f)
        scanner := bufio.NewScanner(file)
        if showNumbers {
            lineCount := 0
            for scanner.Scan() {
                lineCount += 1

                if showEspecialNumbers && scanner.Text() == ""{
                    lineCount -= 1
                    fmt.Fprint(os.Stdout, "\n")
                    continue
                }
                fmt.Fprintf(os.Stdout, "   %d %s\n", lineCount, scanner.Text())
            }
        } else {
            for scanner.Scan() {  
                fmt.Fprintf(os.Stdout, "%s\n", scanner.Text())
            }
        }
    }
}

func isValidFile(file string) error {
    FileInfo, _ := os.Stat(file)

    if FileInfo.IsDir() {
        return errors.New(FileInfo.Name() + " is a directory not a file to cat")
    }
    return nil
}
