// internal/crawler/extractors/herkunft_extractor.go
package extractors

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// HerkunftExtractor extracts word origins
type HerkunftExtractor struct {
	*BaseExtractor
}

// NewHerkunftExtractor creates a new HerkunftExtractor
func NewHerkunftExtractor(doc *goquery.Document) Extractor {
	return &HerkunftExtractor{
		BaseExtractor: NewBaseExtractor(doc),
	}
}

// Extract extracts word origins
func (e *HerkunftExtractor) Extract() interface{} {
	var origins []models.Origin

	// Select the paragraph inside the Herkunft section
	e.Doc.Find("#herkunft p").Each(func(i int, s *goquery.Selection) {
		// Process each element in the paragraph
		s.Contents().Each(func(j int, content *goquery.Selection) {
			if content.Is("a") {
				// If it's a link
				text := e.CleanText(content.Text())
				href, _ := content.Attr("href")

				if text != "" {
					origins = append(origins, models.Origin{
						Word: text,
						Link: href,
					})
				}
			} else {
				// If it's plain text
				text := content.Text()
				if text == "" {
					return
				}

				// Split the text by spaces and add each word
				words := strings.Fields(text)
				for _, word := range words {
					word = e.CleanText(word)
					if word != "" {
						origins = append(origins, models.Origin{
							Word: word,
							Link: "",
						})
					}
				}
			}
		})
	})

	// If no origins found, add a default one
	if len(origins) == 0 {
		origins = append(origins, models.Origin{
			Word: "No data available",
			Link: "",
		})
	}

	return origins
}
