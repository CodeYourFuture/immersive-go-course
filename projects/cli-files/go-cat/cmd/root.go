package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func Execute() {
    if len(os.Args) == 1 {
        fmt.Fprintln(os.Stderr, "Please type a valid file to cat")
        return
    }
    
    if os.Args[1] == "-n" {
        iterateOverArgs(2, true, false)
        return
    }
    if os.Args[1] == "-b" {
        iterateOverArgs(2, true, true)
        return
    }
    iterateOverArgs(1, false, false) 
    return
}    

func iterateOverArgs(startingIndex int, showNumbers bool, showEspecialNumbers bool) {
    for _, f := range os.Args[startingIndex:] {
        if !isValidFile(f) {
            fmt.Fprintln(os.Stderr, "Please type a valid file to cat")
            continue
        }
        file, _ := os.Open(f)
        scanner := bufio.NewScanner(file)
        if showNumbers {
            lineCount := 0
            for scanner.Scan() {
                lineCount += 1

                if showEspecialNumbers {
                    if scanner.Text() == "" {
                        lineCount -= 1
                        fmt.Fprint(os.Stdout, "\n")
                        continue
                    }
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

func isValidFile(file string) bool {
    FileInfo, _ := os.Stat(file)

    if FileInfo.IsDir() {
        return false
    }
    return true
}
