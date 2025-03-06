// File: cmd/test_db.go

package cmd

import (
	"fmt"

	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/container"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/spf13/cobra"
)

// testDbCmd represents the test-db command
var testDbCmd = &cobra.Command{
	Use:   "test-db",
	Short: "Test database connections",
	Long:  `Attempts to connect to all configured databases and reports their status.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Testing database connections...")

		// Get the word repository from the container
		wordRepo := container.GetWordRepository()

		// Try to get a test word from the database
		testWord := "test"
		word, err := wordRepo.GetWord(testWord)

		if err != nil {
			fmt.Printf("Database test failed: %v\n", err)
			fmt.Println("Note: This failure is expected if no data exists yet.")
			return
		}

		fmt.Println("Database test passed!")
		fmt.Printf("Retrieved word: %s\n", word.Word)
	},
}

func init() {
	rootCmd.AddCommand(testDbCmd)
}
