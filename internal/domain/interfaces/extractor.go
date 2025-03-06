// File: internal/domain/interfaces/extractor.go

package interfaces

import (
	"github.com/PuerkitoBio/goquery"
)

// Extractor defines the interface for extracting data from HTML documents
type Extractor interface {
	// Extract extracts data from the document
	Extract() (interface{}, error)

	// GetName returns the name of the extractor
	GetName() string
}

// ExtractorFactory defines the interface for creating extractors
type ExtractorFactory interface {
	// CreateExtractor creates an extractor for the given section
	CreateExtractor(section string, doc *goquery.Document) (Extractor, error)

	// GetAvailableSections returns all available section names
	GetAvailableSections() []string
}

// ExtractorRegistry defines the interface for registering and retrieving extractors
type ExtractorRegistry interface {
	// RegisterExtractor registers an extractor
	RegisterExtractor(name string, constructor func(*goquery.Document) Extractor)

	// GetExtractor retrieves an extractor by name
	GetExtractor(name string) (func(*goquery.Document) Extractor, bool)

	// GetAllExtractors returns all registered extractors
	GetAllExtractors() map[string]func(*goquery.Document) Extractor
}

// BaseExtractorBehavior defines the common behavior for all extractors
type BaseExtractorBehavior interface {
	// ExtractText extracts text from an element matching the selector
	ExtractText(selector string, defaultValue string) string

	// ExtractList extracts a list of texts from elements matching the selector
	ExtractList(selector string) []string
}
