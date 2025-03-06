// File: internal/domain/interfaces/service.go

package interfaces

import (
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// WordDataProvider defines the interface for providing word data
type WordDataProvider interface {
	// GetWordData fetches structured data for a word
	GetWordData(word string) (*models.Word, error)
}

// SuggestionService defines the interface for providing word suggestions
type SuggestionService interface {
	// GetWordSuggestions fetches suggestions for a word
	GetWordSuggestions(word string) ([]models.Synonym, error)
}

// SectionService defines the interface for providing available sections
type SectionService interface {
	// GetAvailableSections returns all available sections that can be extracted
	GetAvailableSections() []string
}

// WordService combines all word-related service interfaces
type WordService interface {
	WordDataProvider
	SuggestionService
	SectionService
}

// FormatterService defines the interface for formatting output
type FormatterService interface {
	// FormatOutput formats the output based on user selection
	FormatOutput(wordData *models.Word, format string) (string, error)
}
