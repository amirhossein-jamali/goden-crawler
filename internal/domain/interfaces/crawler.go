// File: internal/domain/interfaces/crawler.go

package interfaces

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// HTMLFetcher defines the interface for fetching HTML documents
type HTMLFetcher interface {
	// FetchHTML fetches an HTML document from a URL
	FetchHTML(url string) (*goquery.Document, error)
}

// WordFetcher defines the interface for fetching word data
type WordFetcher interface {
	// FetchWord fetches data for a word
	FetchWord(word string) (*models.Word, error)
}

// SuggestionProvider defines the interface for providing word suggestions
type SuggestionProvider interface {
	// GetSuggestions returns a list of suggested words for a given input
	GetSuggestions(word string) ([]models.Synonym, error)
}

// SectionProvider defines the interface for providing available sections
type SectionProvider interface {
	// GetAvailableSections returns a list of all available sections that can be extracted
	GetAvailableSections() []string
}

// Crawler combines all crawler interfaces
type Crawler interface {
	HTMLFetcher
	WordFetcher
	SuggestionProvider
	SectionProvider
}
