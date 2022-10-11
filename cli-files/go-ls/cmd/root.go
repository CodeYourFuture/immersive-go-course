package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRoodCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "go-ls",
		Short: "go-ls is a re-implementation of the ls command",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// By default, ls the current directory
			dir := "."
			// If an argument was passed, use the first as our directory
			if len(args) > 0 {
				dir = args[0]
			}

			// Stat the file so we can check if it's a directory or not before
			// we try to read it as a directory using ReadDir. os.ReadDir will
			// generates an error if the thing you pass to it is not a directory.
			// https://pkg.go.dev/os#Stat
			fileInfo, err := os.Stat(dir)
			if err != nil {
				return err
			}

			// We can only list the contents of a directory.
			// To match the real ls, if we're asked to ls a file, we'll just print
			// out the file's name.
			// https://pkg.go.dev/io/fs#FileInfo
			if fileInfo.IsDir() == false {
				fmt.Fprintln(cmd.OutOrStdout(), fileInfo.Name())
				return nil
			}

			// Read this directory to get a list of files
			// https://pkg.go.dev/os#ReadDir
			files, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			// Iterate through each file in the directory, printing the file name
			for _, file := range files {
				fmt.Fprintln(cmd.OutOrStdout(), file.Name())
			}
			return nil
		},
	}
}

func Execute() {
	rootCmd := NewRoodCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
