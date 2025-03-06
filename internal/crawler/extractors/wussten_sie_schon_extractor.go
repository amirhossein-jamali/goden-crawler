// internal/crawler/extractors/wussten_sie_schon_extractor.go
package extractors

import (
	"github.com/PuerkitoBio/goquery"
)

// WusstenSieSchonExtractor extracts "Did you know?" information
type WusstenSieSchonExtractor struct {
	*BaseExtractor
}

// NewWusstenSieSchonExtractor creates a new WusstenSieSchonExtractor
func NewWusstenSieSchonExtractor(doc *goquery.Document) Extractor {
	return &WusstenSieSchonExtractor{
		BaseExtractor: NewBaseExtractor(doc),
	}
}

// Extract extracts "Did you know?" information
func (e *WusstenSieSchonExtractor) Extract() interface{} {
	var funFacts []string

	e.Doc.Find("#wussten_sie_schon ul li").Each(func(i int, s *goquery.Selection) {
		text := e.CleanText(s.Text())
		if text != "" {
			funFacts = append(funFacts, text)
		}
	})

	return funFacts
}
