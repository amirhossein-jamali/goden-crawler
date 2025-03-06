// File: pkg/logger/logger.go

package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// DEBUG level for detailed debugging information
	DEBUG LogLevel = iota
	// INFO level for general operational information
	INFO
	// WARN level for warning messages
	WARN
	// ERROR level for error messages
	ERROR
	// FATAL level for critical errors that require immediate attention
	FATAL
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns the ANSI color code for a log level
func (l LogLevel) Color() string {
	switch l {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case FATAL:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new log field
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Logger provides logging functionality
type Logger struct {
	level      LogLevel
	output     io.Writer
	prefix     string
	color      bool
	jsonFormat bool
}

// NewLogger creates a new logger with the specified settings
func NewLogger(level LogLevel, output io.Writer, prefix string, color bool) *Logger {
	return &Logger{
		level:      level,
		output:     output,
		prefix:     prefix,
		color:      color,
		jsonFormat: false,
	}
}

// DefaultLogger returns a default logger
func DefaultLogger() *Logger {
	return &Logger{
		level:      INFO,
		output:     os.Stdout,
		prefix:     "",
		color:      true,
		jsonFormat: false,
	}
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

// SetPrefix sets the log prefix
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// SetColor enables or disables colored output
func (l *Logger) SetColor(color bool) {
	l.color = color
}

// SetJSONFormat enables or disables JSON formatted logs
func (l *Logger) SetJSONFormat(jsonFormat bool) {
	l.jsonFormat = jsonFormat
}

// log logs a message with the specified level and fields
func (l *Logger) log(level LogLevel, message string, fields ...Field) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	if l.jsonFormat {
		// Create a map for JSON logging
		logEntry := map[string]interface{}{
			"timestamp": timestamp,
			"level":     level.String(),
			"message":   message,
		}

		if l.prefix != "" {
			logEntry["prefix"] = l.prefix
		}

		// Add fields to the log entry
		for _, field := range fields {
			logEntry[field.Key] = field.Value
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(logEntry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling log entry to JSON: %v\n", err)
			return
		}

		fmt.Fprintln(l.output, string(jsonData))
	} else {
		// Text format logging
		var logLine string

		if l.color {
			levelColor := level.Color()
			resetColor := "\033[0m"

			if l.prefix != "" {
				logLine = fmt.Sprintf("%s [%s%s%s] [%s] %s",
					timestamp, levelColor, level.String(), resetColor, l.prefix, message)
			} else {
				logLine = fmt.Sprintf("%s [%s%s%s] %s",
					timestamp, levelColor, level.String(), resetColor, message)
			}
		} else {
			if l.prefix != "" {
				logLine = fmt.Sprintf("%s [%s] [%s] %s",
					timestamp, level.String(), l.prefix, message)
			} else {
				logLine = fmt.Sprintf("%s [%s] %s",
					timestamp, level.String(), message)
			}
		}

		// Add fields if present
		if len(fields) > 0 {
			logLine += " {"
			for i, field := range fields {
				if i > 0 {
					logLine += ", "
				}
				logLine += fmt.Sprintf("%s=%v", field.Key, field.Value)
			}
			logLine += "}"
		}

		fmt.Fprintln(l.output, logLine)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...Field) {
	l.log(DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...Field) {
	l.log(INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...Field) {
	l.log(WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...Field) {
	l.log(ERROR, message, fields...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(message string, fields ...Field) {
	l.log(FATAL, message, fields...)
	os.Exit(1)
}

// Global logger instance
var globalLogger = DefaultLogger()

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return globalLogger
}

// Debug logs a debug message using the global logger
func Debug(message string, fields ...Field) {
	globalLogger.Debug(message, fields...)
}

// Info logs an info message using the global logger
func Info(message string, fields ...Field) {
	globalLogger.Info(message, fields...)
}

// Warn logs a warning message using the global logger
func Warn(message string, fields ...Field) {
	globalLogger.Warn(message, fields...)
}

// Error logs an error message using the global logger
func Error(message string, fields ...Field) {
	globalLogger.Error(message, fields...)
}

// Fatal logs a fatal message using the global logger and exits the program
func Fatal(message string, fields ...Field) {
	globalLogger.Fatal(message, fields...)
}

// Deprecated functions for backward compatibility

// Debugf logs a formatted debug message (deprecated, use Debug with fields)
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message (deprecated, use Info with fields)
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message (deprecated, use Warn with fields)
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message (deprecated, use Error with fields)
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Fatalf logs a formatted fatal message and exits (deprecated, use Fatal with fields)
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// Debugf logs a formatted debug message using the global logger (deprecated)
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Infof logs a formatted info message using the global logger (deprecated)
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Warnf logs a formatted warning message using the global logger (deprecated)
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Errorf logs a formatted error message using the global logger (deprecated)
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message using the global logger and exits (deprecated)
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}
