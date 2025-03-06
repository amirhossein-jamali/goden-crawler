// File: pkg/utils/config.go

package utils

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	// HTTP settings
	HTTPTimeout time.Duration
	HTTPRetries int
	UserAgent   string

	// Duden settings
	DudenBaseURL   string
	DudenSearchURL string

	// Logging settings
	LogLevel        string
	EnableColorLogs bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		HTTPTimeout:     10 * time.Second,
		HTTPRetries:     3,
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		DudenBaseURL:    "https://www.duden.de",
		DudenSearchURL:  "https://www.duden.de/suchen/dudenonline/",
		LogLevel:        "INFO",
		EnableColorLogs: true,
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := DefaultConfig()

	// Load HTTP settings
	if timeout, err := strconv.Atoi(getEnv("HTTP_TIMEOUT_SECONDS", "10")); err == nil {
		config.HTTPTimeout = time.Duration(timeout) * time.Second
	}

	if retries, err := strconv.Atoi(getEnv("HTTP_RETRIES", "3")); err == nil {
		config.HTTPRetries = retries
	}

	if userAgent := getEnv("USER_AGENT", ""); userAgent != "" {
		config.UserAgent = userAgent
	}

	// Load Duden settings
	if baseURL := getEnv("DUDEN_BASE_URL", ""); baseURL != "" {
		config.DudenBaseURL = baseURL
	}

	if searchURL := getEnv("DUDEN_SEARCH_URL", ""); searchURL != "" {
		config.DudenSearchURL = searchURL
	}

	// Load logging settings
	if logLevel := getEnv("LOG_LEVEL", ""); logLevel != "" {
		config.LogLevel = logLevel
	}

	if colorLogs, err := strconv.ParseBool(getEnv("ENABLE_COLOR_LOGS", "true")); err == nil {
		config.EnableColorLogs = colorLogs
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
