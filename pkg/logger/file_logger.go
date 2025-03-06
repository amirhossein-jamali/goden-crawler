// File: pkg/logger/file_logger.go

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var fileLogger *log.Logger
var logFile *os.File

// InitFileLogger initializes a logger that writes to a file
func InitFileLogger() {
	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join("logs", fmt.Sprintf("goden-crawler_%s.log", timestamp))

	logFile, err = os.Create(logFilePath)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	fileLogger = log.New(logFile, "", log.LstdFlags)

	Info("File logger initialized", F("path", logFilePath))
}

// LogToFile logs a message to the file
func LogToFile(message string) {
	if fileLogger == nil {
		InitFileLogger()
	}

	fileLogger.Println(message)
}
