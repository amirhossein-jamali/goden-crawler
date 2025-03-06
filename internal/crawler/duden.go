// internal/crawler/duden.go
package crawler

import (
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// DudenCrawler defines the interface for scraping data from the Duden website.
// This interface allows for different implementations and makes testing easier.
type DudenCrawler interface {
	// FetchWordData fetches data for a word and returns a map of section -> data
	// This is useful for getting raw data without processing it into structured models.
	FetchWordData(word string) (map[string]string, error)

	// FetchWordDataStructured fetches data for a word and returns a structured Word object
	// This is the primary method for getting complete, structured linguistic data.
	FetchWordDataStructured(word string) (*models.Word, error)

	// GetSuggestions returns a list of suggested words for a given input
	// Useful when the exact word is not found but similar words exist.
	GetSuggestions(word string) ([]models.Synonym, error)

	// GetAvailableSections returns a list of all available sections that can be extracted
	// This helps clients know what data is available for extraction.
	GetAvailableSections() []string
}
