// File: internal/crawler/extractors/base_extractor.go

package extractors

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Extractor is the interface that all extractors must implement
type Extractor interface {
	Extract() interface{}
}

// BaseExtractor provides common functionality for all extractors
type BaseExtractor struct {
	Doc *goquery.Document
}

// NewBaseExtractor creates a new BaseExtractor
func NewBaseExtractor(doc *goquery.Document) *BaseExtractor {
	return &BaseExtractor{
		Doc: doc,
	}
}

// CleanText removes hidden characters and unnecessary symbols
func (b *BaseExtractor) CleanText(text string) string {
	if text == "" {
		return ""
	}

	// Remove soft hyphens and other special characters
	text = strings.ReplaceAll(text, "\u00ad", "")
	text = strings.ReplaceAll(text, "ⓘ", "")
	text = strings.ReplaceAll(text, "→", "")

	return strings.TrimSpace(text)
}

// ExtractText extracts text from an element matching the selector
func (b *BaseExtractor) ExtractText(selector string, defaultValue string) string {
	var text string
	b.Doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text = s.Text()
	})

	if text == "" {
		return defaultValue
	}

	return b.CleanText(text)
}

// ExtractList extracts a list of texts from elements matching the selector
func (b *BaseExtractor) ExtractList(selector string) []string {
	var results []string
	b.Doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := b.CleanText(s.Text())
		if text != "" {
			results = append(results, text)
		}
	})
	return results
}

// registry keeps track of all registered extractors
var registry = make(map[string]Extractor)

// RegisterExtractor automatically registers an extractor
func RegisterExtractor(name string, extractor Extractor) {
	registry[name] = extractor
}

// GetExtractor retrieves an extractor by name
func GetExtractor(name string) (Extractor, bool) {
	extractor, exists := registry[name]
	return extractor, exists
}

// GetAllExtractors returns a list of all registered extractors
func GetAllExtractors() map[string]Extractor {
	return registry
}
