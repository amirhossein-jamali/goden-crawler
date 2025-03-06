// File: cmd/interactive.go

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/amirhossein-jamali/goden-crawler/internal/application/services"
	"github.com/amirhossein-jamali/goden-crawler/internal/formatter"
	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/container"

	"github.com/spf13/cobra"
)

var interactiveFormat = "text"

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start an interactive session",
	Long: `Start an interactive session where you can enter words to scrape.
Type 'exit' or 'quit' to end the session.
Type 'help' to see available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get word service from container
		wordService := container.GetWordService()

		// Start interactive mode
		fmt.Println("🔍 Goden Crawler Interactive Mode")
		fmt.Println("Type a German word to get information, or 'exit' to quit.")
		fmt.Println("Type 'help' for more commands.")

		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Print("\n> ")
			if !scanner.Scan() {
				break
			}

			input := strings.TrimSpace(scanner.Text())

			if input == "" {
				continue
			}

			// Process commands
			switch strings.ToLower(input) {
			case "exit", "quit":
				fmt.Println("👋 Goodbye!")
				return

			case "help":
				printHelp()
				continue

			case "sections":
				sections := wordService.GetAvailableSections()
				fmt.Println("📋 Available sections:")
				for i, section := range sections {
					fmt.Printf("  %d. %s\n", i+1, section)
				}
				continue

			case "format json":
				interactiveFormat = "json"
				fmt.Println("🔄 Output format set to JSON")
				continue

			case "format text":
				interactiveFormat = "text"
				fmt.Println("🔄 Output format set to text")
				continue
			}

			// Check if it's a suggestion request
			if strings.HasPrefix(input, "suggest ") {
				word := strings.TrimPrefix(input, "suggest ")
				handleSuggestions(wordService, word)
				continue
			}

			// Otherwise, treat as a word to scrape
			handleWord(wordService, input, interactiveFormat)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("🚨 Error reading input:", err)
		}
	},
}

func handleWord(wordService *services.WordService, word, format string) {
	fmt.Printf("🔍 Fetching information for '%s'...\n", word)

	// Fetch word data using the service
	wordData, err := wordService.GetWordData(word)
	if err != nil {
		fmt.Println("🚨 Error:", err)
		return
	}

	// Format output based on user selection
	output, err := formatter.FormatOutput(wordData, format)
	if err != nil {
		fmt.Println("🚨 Error formatting output:", err)
		return
	}

	fmt.Println(output)
}

func handleSuggestions(wordService *services.WordService, word string) {
	fmt.Printf("🔍 Finding suggestions for '%s'...\n", word)

	suggestions, err := wordService.GetWordSuggestions(word)
	if err != nil {
		fmt.Println("🚨 Error:", err)
		return
	}

	if len(suggestions) == 0 {
		fmt.Println("❌ No suggestions found.")
		return
	}

	fmt.Printf("📋 Found %d suggestions:\n", len(suggestions))
	for i, suggestion := range suggestions {
		fmt.Printf("  %d. %s\n", i+1, suggestion.Text)
	}
}

func printHelp() {
	fmt.Println("📚 Available commands:")
	fmt.Println("  [word]           - Fetch information for a German word")
	fmt.Println("  suggest [word]   - Get suggestions for a word")
	fmt.Println("  format json      - Set output format to JSON")
	fmt.Println("  format text      - Set output format to text")
	fmt.Println("  sections         - List available data sections")
	fmt.Println("  help             - Show this help message")
	fmt.Println("  exit, quit       - Exit the program")
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
