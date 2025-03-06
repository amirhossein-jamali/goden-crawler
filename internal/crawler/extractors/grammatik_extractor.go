// internal/crawler/extractors/grammatik_extractor.go
package extractors

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GrammatikExtractor extracts grammatical information
type GrammatikExtractor struct {
	*BaseExtractor
}

// NewGrammatikExtractor creates a new GrammatikExtractor
func NewGrammatikExtractor(doc *goquery.Document) Extractor {
	return &GrammatikExtractor{
		BaseExtractor: NewBaseExtractor(doc),
	}
}

// Extract extracts grammatical information
func (e *GrammatikExtractor) Extract() interface{} {
	return map[string]interface{}{
		"paragraphs": e.extractParagraphs(),
		"links":      e.extractLinks(),
	}
}

// extractLinks extracts all links from the 'Grammatik' section
func (e *GrammatikExtractor) extractLinks() []map[string]string {
	var linksData []map[string]string

	e.Doc.Find("#grammatik a.more__link").Each(func(i int, s *goquery.Selection) {
		text := e.CleanText(s.Text())
		link, exists := s.Attr("href")
		if !exists {
			link = "N/A"
		}

		linksData = append(linksData, map[string]string{
			"text": text,
			"link": link,
		})
	})

	if len(linksData) == 0 {
		linksData = append(linksData, map[string]string{
			"text": "No links available",
			"link": "N/A",
		})
	}

	return linksData
}

// extractParagraphs extracts and parses the grammatical information paragraph
func (e *GrammatikExtractor) extractParagraphs() map[string]interface{} {
	var text string
	var details []map[string]string

	e.Doc.Find("#grammatik p").Each(func(i int, s *goquery.Selection) {
		text = e.CleanText(s.Text())

		// Split the text into parts based on ';'
		parts := strings.Split(text, ";")

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			if strings.Contains(part, ":") {
				// Split by ':' to get key-value pairs
				keyValue := strings.SplitN(part, ":", 2)
				key := strings.TrimSpace(keyValue[0])
				value := strings.TrimSpace(keyValue[1])

				// Split the value by ',' and process sub-values
				subValues := strings.Split(value, ",")
				for _, subValue := range subValues {
					subValue = strings.TrimSpace(subValue)
					if subValue == "" {
						continue
					}

					if strings.Contains(subValue, ":") {
						// If sub-value contains another key-value pair
						subKeyValue := strings.SplitN(subValue, ":", 2)
						subKey := strings.TrimSpace(subKeyValue[0])
						subVal := strings.TrimSpace(subKeyValue[1])

						details = append(details, map[string]string{
							"type":  subKey,
							"value": subVal,
						})
					} else {
						details = append(details, map[string]string{
							"type":  key,
							"value": subValue,
						})
					}
				}
			} else {
				details = append(details, map[string]string{
					"type":  "base_form",
					"value": part,
				})
			}
		}
	})

	if text == "" {
		text = "No data"
	}

	return map[string]interface{}{
		"text":    text,
		"details": details,
	}
}
