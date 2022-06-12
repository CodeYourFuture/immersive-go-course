package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-ls",
	Short: "go-ls is a re-implementation of the ls command",
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := os.ReadDir(".")
		if err != nil {
			return err
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
