// internal/crawler/extractors/synonyme_extractor.go
package extractors

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// SynonymeExtractor extracts synonyms
type SynonymeExtractor struct {
	*BaseExtractor
}

// NewSynonymeExtractor creates a new SynonymeExtractor
func NewSynonymeExtractor(doc *goquery.Document) Extractor {
	return &SynonymeExtractor{
		BaseExtractor: NewBaseExtractor(doc),
	}
}

// Extract extracts synonyms
func (e *SynonymeExtractor) Extract() interface{} {
	return map[string]interface{}{
		"synonyms":  e.extractSynonyms(),
		"more_link": e.extractMoreLink(),
	}
}

// extractSynonyms extracts synonyms from the main page
func (e *SynonymeExtractor) extractSynonyms() []models.Synonym {
	var synonyms []models.Synonym

	e.Doc.Find("#synonyme ul li").Each(func(i int, s *goquery.Selection) {
		// Get the text and split by commas
		text := e.CleanText(s.Text())
		parts := strings.Split(text, ",")

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// Check if this part has a link
			var href string
			s.Find("a").Each(func(j int, a *goquery.Selection) {
				if e.CleanText(a.Text()) == part {
					if h, exists := a.Attr("href"); exists {
						href = h
					}
				}
			})

			synonyms = append(synonyms, models.Synonym{
				Text: part,
				Link: href,
			})
		}
	})

	// If we have a "more link", fetch additional synonyms
	moreLink := e.extractMoreLink()
	if moreLink != "" {
		additionalSynonyms := e.fetchAdditionalSynonyms(moreLink)
		synonyms = append(synonyms, additionalSynonyms...)
	}

	return synonyms
}

// extractMoreLink extracts the "more link" if available
func (e *SynonymeExtractor) extractMoreLink() string {
	var moreLink string
	e.Doc.Find("#synonyme .more__link").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			moreLink = href
		}
	})
	return moreLink
}

// fetchAdditionalSynonyms fetches additional synonyms from the provided link
func (e *SynonymeExtractor) fetchAdditionalSynonyms(link string) []models.Synonym {
	var additionalSynonyms []models.Synonym

	// Ensure the link is absolute
	if !strings.HasPrefix(link, "http") {
		link = "https://www.duden.de" + link
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Get(link)
	if err != nil {
		return additionalSynonyms
	}
	defer resp.Body.Close()

	// Parse the response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return additionalSynonyms
	}

	// Extract synonyms from the additional page
	doc.Find(".content-section .vignette__content").Each(func(i int, s *goquery.Selection) {
		text := e.CleanText(s.Text())
		parts := strings.Split(text, ",")

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			additionalSynonyms = append(additionalSynonyms, models.Synonym{
				Text: part,
				Link: "",
			})
		}
	})

	return additionalSynonyms
}
