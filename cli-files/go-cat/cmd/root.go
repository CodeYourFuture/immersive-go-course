/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
func NewRoodCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "go-cat",
		Short: "Go implementation of cat",
		Long:  `Works like cat`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Don't do anything if we didn't get an arg
			if len(args) < 1 {
				return nil
			}
			path := args[0]

			// Get data about the file so we can do this safely
			file, err := os.Stat(path)
			if err != nil {
				return err
			}

			// If it's a directory, do the right thing and error
			if file.IsDir() {
				return fmt.Errorf("go-cat: %s: Is a directory", path)
			}

			// Read the data from the file
			// https://pkg.go.dev/os#ReadFile
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Print those lovely bytes
			out := cmd.OutOrStdout()
			out.Write(data)

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
