// File: cmd/bulk.go

package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/internal/application/services"
	"github.com/amirhossein-jamali/goden-crawler/internal/formatter"
	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/container"
	"github.com/amirhossein-jamali/goden-crawler/internal/repository"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	bulkInputFile  string
	bulkOutputDir  string
	bulkBatchSize  int
	bulkWorkers    int
	bulkTimeoutSec int
	bulkFormat     string
)

// bulkCmd represents the bulk command
var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Process a large number of words from a file",
	Long: `Process a large number of German words from a file.
Each word should be on a separate line in the input file.
Results will be saved to the specified output directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate input file
		if bulkInputFile == "" {
			fmt.Println("Error: Input file is required")
			cmd.Help()
			os.Exit(1)
		}

		// Check if input file exists
		if _, err := os.Stat(bulkInputFile); os.IsNotExist(err) {
			fmt.Printf("Error: Input file '%s' not found\n", bulkInputFile)
			os.Exit(1)
		}

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(bulkOutputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}

		// Read words from file
		words, err := readWordsFromFile(bulkInputFile)
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}

		if len(words) == 0 {
			fmt.Println("No words found in the input file")
			os.Exit(1)
		}

		fmt.Printf("Found %d words in the input file\n", len(words))
		fmt.Printf("Processing with batch size: %d, workers: %d, timeout: %d seconds\n",
			bulkBatchSize, bulkWorkers, bulkTimeoutSec)

		// Get services from container
		wordService := container.GetWordService()
		wordRepository := container.GetWordRepository()

		// Process words in batches
		totalSuccess := 0
		totalFailure := 0

		for i := 0; i < len(words); i += bulkBatchSize {
			end := i + bulkBatchSize
			if end > len(words) {
				end = len(words)
			}

			batchWords := words[i:end]
			fmt.Printf("\nProcessing batch %d/%d (%d words)...\n",
				(i/bulkBatchSize)+1, (len(words)+bulkBatchSize-1)/bulkBatchSize, len(batchWords))

			success, failure := processBatch(batchWords, wordService, wordRepository)
			totalSuccess += success
			totalFailure += failure

			// Add a small delay between batches to avoid overwhelming the server
			if end < len(words) {
				fmt.Println("Waiting 2 seconds before next batch...")
				time.Sleep(2 * time.Second)
			}
		}

		// Print summary
		fmt.Printf("\nBulk processing complete:\n")
		fmt.Printf("- Total words: %d\n", len(words))
		fmt.Printf("- Successful: %d\n", totalSuccess)
		fmt.Printf("- Failed: %d\n", totalFailure)
		fmt.Printf("Results saved to: %s\n", bulkOutputDir)

		if totalFailure > 0 {
			os.Exit(1)
		}
	},
}

// readWordsFromFile reads words from a file, one word per line
func readWordsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" && !strings.HasPrefix(word, "#") {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// processBatch processes a batch of words
func processBatch(words []string, wordService *services.WordService, wordRepository *repository.WordRepository) (int, int) {
	// Create a batch service
	batchService := services.NewBatchService(
		wordService,
		wordRepository,
		bulkWorkers,
		time.Duration(bulkTimeoutSec)*time.Second,
	)

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

		// Save the output to a file
		filename := filepath.Join(bulkOutputDir, fmt.Sprintf("%s.%s", result.Word, bulkFormat))

		// Format the output
		output, err := formatter.FormatOutput(result.Data, bulkFormat)
		if err != nil {
			logger.Error("Failed to format output",
				logger.F("word", result.Word),
				logger.F("error", err))
			failureCount++
			continue
		}

		// Write to file
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

	return successCount, failureCount
}

func init() {
	rootCmd.AddCommand(bulkCmd)

	// Add flags
	bulkCmd.Flags().StringVarP(&bulkInputFile, "input", "i", "", "Input file containing words (required)")
	bulkCmd.Flags().StringVarP(&bulkOutputDir, "output", "o", "output", "Output directory for results")
	bulkCmd.Flags().IntVarP(&bulkBatchSize, "batch-size", "b", 10, "Number of words to process in each batch")
	bulkCmd.Flags().IntVarP(&bulkWorkers, "workers", "w", 5, "Number of concurrent workers")
	bulkCmd.Flags().IntVarP(&bulkTimeoutSec, "timeout", "t", 30, "Timeout in seconds per word")
	bulkCmd.Flags().StringVarP(&bulkFormat, "format", "f", "json", "Output format (text, json)")

	// Mark required flags
	bulkCmd.MarkFlagRequired("input")
}
