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
