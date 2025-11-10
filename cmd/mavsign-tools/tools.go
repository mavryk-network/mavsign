package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mavsign-tools",
		Short: "Various MavSign tools",
	}

	rootCmd.AddCommand(NewGenKeyCommand())
	rootCmd.AddCommand(NewAuthRequestCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
