package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fffetch",
	Short: "Fantasy Football Data Fetcher",
	Long:  "Fetch and process fantasy football data from Pro Football Reference",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
