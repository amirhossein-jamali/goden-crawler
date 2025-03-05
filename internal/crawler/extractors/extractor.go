// ./internal/crawler/extractors/extractor.go

package extractors

import "github.com/PuerkitoBio/goquery"

// Extractor defines the interface for all extractors
type Extractor interface {
	Extract(doc *goquery.Document) string
}
