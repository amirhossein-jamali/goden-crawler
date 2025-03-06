// File: pkg/errors/errors.go

package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Standard errors that can be used throughout the application
var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrNetworkError  = errors.New("network error")
	ErrParsingError  = errors.New("parsing error")
	ErrInternalError = errors.New("internal error")
)

// DudenScraperError represents an error that occurred during scraping
type DudenScraperError struct {
	Word    string
	Message string
	Err     error
	Stack   string
}

// Error implements the error interface
func (e *DudenScraperError) Error() string {
	if e.Word != "" {
		return fmt.Sprintf("error scraping word '%s': %s", e.Word, e.Message)
	}
	return fmt.Sprintf("scraper error: %s", e.Message)
}

// Unwrap returns the underlying error
func (e *DudenScraperError) Unwrap() error {
	return e.Err
}

// NewScraperError creates a new DudenScraperError
func NewScraperError(word, message string, err error) *DudenScraperError {
	return &DudenScraperError{
		Word:    word,
		Message: message,
		Err:     err,
		Stack:   captureStack(2),
	}
}

// ExtractorError represents an error that occurred during extraction
type ExtractorError struct {
	Section string
	Message string
	Err     error
	Stack   string
}

// Error implements the error interface
func (e *ExtractorError) Error() string {
	return fmt.Sprintf("error extracting section '%s': %s", e.Section, e.Message)
}

// Unwrap returns the underlying error
func (e *ExtractorError) Unwrap() error {
	return e.Err
}

// NewExtractorError creates a new ExtractorError
func NewExtractorError(section, message string, err error) *ExtractorError {
	return &ExtractorError{
		Section: section,
		Message: message,
		Err:     err,
		Stack:   captureStack(2),
	}
}

// Wrap wraps an error with a message and returns a new error
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapWithContext wraps an error with context information
func WrapWithContext(err error, message string, context map[string]interface{}) error {
	if err == nil {
		return nil
	}

	contextStr := formatContext(context)
	if contextStr != "" {
		return fmt.Errorf("%s [%s]: %w", message, contextStr, err)
	}

	return fmt.Errorf("%s: %w", message, err)
}

// formatContext formats a context map as a string
func formatContext(context map[string]interface{}) string {
	if len(context) == 0 {
		return ""
	}

	var parts []string
	for k, v := range context {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(parts, ", ")
}

// captureStack captures the current stack trace
func captureStack(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		// Skip runtime and standard library frames
		if strings.Contains(frame.File, "runtime/") {
			continue
		}

		fmt.Fprintf(&builder, "%s:%d %s\n", frame.File, frame.Line, frame.Function)

		// Limit stack trace to a reasonable depth
		if builder.Len() > 1000 {
			builder.WriteString("...\n")
			break
		}
	}

	return builder.String()
}
