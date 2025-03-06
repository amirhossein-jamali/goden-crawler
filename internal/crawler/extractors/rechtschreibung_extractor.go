// internal/crawler/extractors/rechtschreibung_extractor.go
package extractors

import (
	"github.com/PuerkitoBio/goquery"
)

// RechtschreibungExtractor extracts spelling information
type RechtschreibungExtractor struct {
	*BaseExtractor
}

// NewRechtschreibungExtractor creates a new RechtschreibungExtractor
func NewRechtschreibungExtractor(doc *goquery.Document) Extractor {
	return &RechtschreibungExtractor{
		BaseExtractor: NewBaseExtractor(doc),
	}
}

// Extract extracts spelling information
func (e *RechtschreibungExtractor) Extract() interface{} {
	return map[string]interface{}{
		"spelling": e.extractSpelling(),
	}
}

// extractSpelling extracts spelling-related information
func (e *RechtschreibungExtractor) extractSpelling() map[string]interface{} {
	// Extract syllabic division (Worttrennung)
	syllabicDivision := e.ExtractText(
		"#rechtschreibung .tuple__key:contains('Worttrennung') + .tuple__val",
		"N/A",
	)

	// Extract general spelling examples from the first list
	generalExamples := e.ExtractList("#rechtschreibung .infobox > ul.infobox__examples > li")

	// Extract examples related to grammar rules (second list)
	ruleRelatedExamples := e.ExtractList("#rechtschreibung .infobox > p + ul.infobox__examples > li")

	// Combine both example lists
	allExamples := append(generalExamples, ruleRelatedExamples...)

	// Extract links to grammatical rules
	var rules []map[string]string
	e.Doc.Find("#rechtschreibung .infobox p a.rule-ref").Each(func(i int, s *goquery.Selection) {
		text := e.CleanText(s.Text())
		href, exists := s.Attr("href")
		if !exists {
			href = ""
		}

		rules = append(rules, map[string]string{
			"text": text,
			"link": href,
		})
	})

	return map[string]interface{}{
		"syllabic_division": syllabicDivision,
		"examples":          allExamples,
		"rules":             rules,
	}
}
