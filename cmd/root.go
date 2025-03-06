// File: cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/container"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goden-crawler",
	Short: "A CLI tool to extract linguistic information from Duden",
	Long: `Goden Crawler is a command-line interface (CLI) tool designed to extract 
detailed linguistic information from the Duden online dictionary.

Built with Golang and Cobra for CLI management, it features a modular 
and scalable architecture, making it easy to maintain and extend.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health of database connections",
	Long:  `Checks if the application can connect to all configured databases.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting database health check...")

		// Get repository from container
		repo := container.GetWordRepository()

		// Try to retrieve a test word
		word, err := repo.GetWord("test")
		if err != nil {
			fmt.Println("Database health check failed, but this is expected if no data exists yet.")
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Database health check passed!")
			fmt.Printf("Retrieved word: %s\n", word.Word)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize logger
	logger.SetGlobalLogger(logger.DefaultLogger())

	// Add commands
	rootCmd.AddCommand(scrapeCmd)
	rootCmd.AddCommand(batchCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(healthCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goden-crawler.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
