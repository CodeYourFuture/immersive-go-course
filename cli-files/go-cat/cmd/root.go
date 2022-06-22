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
		Run: func(cmd *cobra.Command, args []string) {
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
