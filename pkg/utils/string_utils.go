// File: pkg/utils/string_utils.go

package utils

import (
	"strings"
	"unicode"
)

// CleanText removes hidden characters and unnecessary symbols
func CleanText(text string) string {
	if text == "" {
		return ""
	}

	// Remove soft hyphens and other special characters
	text = strings.ReplaceAll(text, "\u00ad", "")
	text = strings.ReplaceAll(text, "ⓘ", "")
	text = strings.ReplaceAll(text, "→", "")

	return strings.TrimSpace(text)
}

// SplitAndTrim splits a string by a separator and trims each part
func SplitAndTrim(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// ToTitleCase converts a string to TitleCase
func ToTitleCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 0; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
