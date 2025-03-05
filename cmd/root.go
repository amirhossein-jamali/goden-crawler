// ./cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Import path should match the module name
var rootCmd = &cobra.Command{
	Use:   "goden-crawler",
	Short: "Goden Crawler is a CLI tool to scrape linguistic data from Duden",
	Long:  `Goden Crawler scrapes detailed linguistic information from the Duden online dictionary.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Goden Crawler! Use 'scrape <word>' to fetch data from Duden.")
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
