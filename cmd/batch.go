// File: cmd/batch.go

package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/internal/application/services"
	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
	"github.com/amirhossein-jamali/goden-crawler/internal/formatter"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	batchFormat  string
	workers      int
	timeoutSecs  int
	outputPrefix string
)

var batchCmd = &cobra.Command{
	Use:   "batch [words]",
	Short: "Process multiple words concurrently",
	Long: `Process multiple German words concurrently and save the results.
Words should be provided as a space-separated list.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		words := args

		// Create a new scraper
		scraper := crawler.NewDudenScraper()

		// Create a batch service
		batchService := services.NewBatchService(
			scraper,
			workers,
			time.Duration(timeoutSecs)*time.Second,
		)

		logger.Info("Starting batch processing",
			logger.F("words", len(words)),
			logger.F("workers", workers),
			logger.F("timeout", timeoutSecs))

		// Process words concurrently
		ctx := context.Background()
		results := batchService.ProcessWords(ctx, words)

		// Process results
		successCount := 0
		failureCount := 0

		for _, result := range results {
			if result.Error != nil {
				logger.Error("Failed to process word",
					logger.F("word", result.Word),
					logger.F("error", result.Error))
				failureCount++
				continue
			}

			// Format the output
			output, err := formatter.FormatOutput(result.Data, batchFormat)
			if err != nil {
				logger.Error("Failed to format output",
					logger.F("word", result.Word),
					logger.F("error", err))
				failureCount++
				continue
			}

			// Save the output to a file
			filename := fmt.Sprintf("%s%s.%s", outputPrefix, result.Word, batchFormat)
			err = os.WriteFile(filename, []byte(output), 0644)
			if err != nil {
				logger.Error("Failed to save output",
					logger.F("word", result.Word),
					logger.F("filename", filename),
					logger.F("error", err))
				failureCount++
				continue
			}

			logger.Info("Successfully processed word",
				logger.F("word", result.Word),
				logger.F("filename", filename))
			successCount++
		}

		// Print summary
		fmt.Printf("\nBatch processing complete:\n")
		fmt.Printf("- Total words: %d\n", len(words))
		fmt.Printf("- Successful: %d\n", successCount)
		fmt.Printf("- Failed: %d\n", failureCount)

		if failureCount > 0 {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringVarP(&batchFormat, "output", "o", "json", "Output format (text, json)")
	batchCmd.Flags().IntVarP(&workers, "workers", "w", 5, "Number of concurrent workers")
	batchCmd.Flags().IntVarP(&timeoutSecs, "timeout", "t", 30, "Timeout in seconds per word")
	batchCmd.Flags().StringVarP(&outputPrefix, "prefix", "p", "", "Output filename prefix")
}
