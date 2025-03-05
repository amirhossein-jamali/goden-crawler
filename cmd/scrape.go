// ./cmd/scrape.go
package cmd

import (
	"fmt"
	"os"

	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
	"github.com/amirhossein-jamali/goden-crawler/internal/formatter"

	"github.com/spf13/cobra"
)

var format string

var scrapeCmd = &cobra.Command{
	Use:   "scrape [word]",
	Short: "Scrapes linguistic data from Duden for a given word",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		word := args[0]

		// Fetch word data using the crawler
		wordData, err := crawler.CrawlDuden(word)
		if err != nil {
			fmt.Println("🚨 Error fetching word data:", err)
			os.Exit(1)
		}

		// Format output based on user selection
		output, err := formatter.FormatOutput(wordData, format)
		if err != nil {
			fmt.Println("🚨 Error formatting output:", err)
			os.Exit(1)
		}

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.Flags().StringVarP(&format, "output", "o", "text", "Output format (text, json)")
}
