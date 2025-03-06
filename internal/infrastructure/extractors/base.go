// File: internal/infrastructure/extractors/base.go

package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/utils"
)

// Extractor is the interface that all extractors must implement
type Extractor interface {
	Extract() interface{}
	GetName() string
}

// BaseExtractor provides common functionality for all extractors
type BaseExtractor struct {
	Doc  *goquery.Document
	Name string
}

// NewBaseExtractor creates a new BaseExtractor
func NewBaseExtractor(doc *goquery.Document, name string) *BaseExtractor {
	return &BaseExtractor{
		Doc:  doc,
		Name: name,
	}
}

// GetName returns the name of the extractor
func (b *BaseExtractor) GetName() string {
	return b.Name
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

	return utils.CleanText(text)
}

// ExtractList extracts a list of texts from elements matching the selector
func (b *BaseExtractor) ExtractList(selector string) []string {
	var results []string
	b.Doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := utils.CleanText(s.Text())
		if text != "" {
			results = append(results, text)
		}
	})
	return results
}

// registry keeps track of all registered extractors
var registry = make(map[string]func(*goquery.Document) Extractor)

// RegisterExtractor registers an extractor constructor function
func RegisterExtractor(name string, constructor func(*goquery.Document) Extractor) {
	registry[name] = constructor
}

// RegisterExtractorType is a helper function to register an extractor type
// This makes registration more similar to Python's decorator pattern
// Usage:
//
//	func init() {
//	    RegisterExtractorType[MyExtractor]("MyExtractor")
//	}
func RegisterExtractorType[T Extractor](name string, constructorFn func(*goquery.Document) T) {
	RegisterExtractor(name, func(doc *goquery.Document) Extractor {
		return constructorFn(doc)
	})
}

// GetExtractor retrieves an extractor constructor by name
func GetExtractor(name string) (func(*goquery.Document) Extractor, bool) {
	constructor, exists := registry[name]
	return constructor, exists
}

// GetAllExtractors returns all registered extractor constructors
func GetAllExtractors() map[string]func(*goquery.Document) Extractor {
	return registry
}
