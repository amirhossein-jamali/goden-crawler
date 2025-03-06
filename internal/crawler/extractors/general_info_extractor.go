// internal/crawler/extractors/general_info_extractor.go
package extractors

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// GeneralInfoExtractor extracts general information about a word
type GeneralInfoExtractor struct {
	*BaseExtractor
}

// NewGeneralInfoExtractor creates a new GeneralInfoExtractor
func NewGeneralInfoExtractor(doc *goquery.Document) Extractor {
	return &GeneralInfoExtractor{
		BaseExtractor: NewBaseExtractor(doc),
	}
}

// Extract extracts general information from the page
func (e *GeneralInfoExtractor) Extract() interface{} {
	return map[string]interface{}{
		"word":          e.extractWord(),
		"article":       e.extractArticle(),
		"word_type":     e.extractWordType(),
		"frequency":     e.extractFrequency(),
		"pronunciation": e.extractPronunciation(),
	}
}

// extractWord extracts the word from the page
func (e *GeneralInfoExtractor) extractWord() string {
	return e.ExtractText(".lemma__main", "")
}

// extractArticle extracts the article from the page
func (e *GeneralInfoExtractor) extractArticle() string {
	return e.ExtractText(".lemma__determiner", "")
}

// extractWordType extracts the word type from the page
func (e *GeneralInfoExtractor) extractWordType() []string {
	wordTypeText := e.ExtractText(".tuple__key:-soup-contains('Wortart') + .tuple__val", "")
	if wordTypeText == "" {
		return []string{"unknown"}
	}

	// Split by comma and trim spaces
	wordTypes := strings.Split(wordTypeText, ",")
	for i, wt := range wordTypes {
		wordTypes[i] = strings.TrimSpace(wt)
	}

	return wordTypes
}

// extractFrequency extracts the frequency rating
func (e *GeneralInfoExtractor) extractFrequency() string {
	var frequency string
	e.Doc.Find(".tuple__key:contains('Häufigkeit') + .tuple__val .shaft").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		filledBars := strings.Count(text, "▒")

		switch filledBars {
		case 5:
			frequency = "very_high"
		case 4:
			frequency = "high"
		case 3:
			frequency = "medium"
		case 2:
			frequency = "low"
		case 1:
			frequency = "very_low"
		default:
			frequency = "unknown"
		}
	})

	return frequency
}

// extractPronunciation extracts pronunciation information
func (e *GeneralInfoExtractor) extractPronunciation() []models.Pronunciation {
	var pronunciations []models.Pronunciation

	e.Doc.Find(".pronunciation-guide").Each(func(i int, s *goquery.Selection) {
		// Extract phonetic transcription
		phonetic := s.Find(".ipa").Text()
		phonetic = e.CleanText(phonetic)

		// Extract audio link
		audioLink := "no_audio_available"
		s.Find("a.pronunciation-guide__sound[data-duden-ref-type='audio']").Each(func(i int, a *goquery.Selection) {
			if href, exists := a.Attr("href"); exists {
				audioLink = href
			}
		})

		// Extract word variants with stress patterns
		s.Find(".pronunciation-guide__text").Each(func(i int, w *goquery.Selection) {
			formattedWord := ""

			// Process each part of the word
			w.Contents().Each(func(i int, c *goquery.Selection) {
				if c.Is(".short-stress") {
					formattedWord += "(" + c.Text() + ")"
				} else if c.Is(".long-stress") {
					formattedWord += "{" + c.Text() + "}"
				} else {
					formattedWord += c.Text()
				}
			})

			formattedWord = e.CleanText(formattedWord)

			if formattedWord != "" {
				pronunciations = append(pronunciations, models.Pronunciation{
					Word:     formattedWord,
					Phonetic: phonetic,
					Audio:    audioLink,
				})
			}
		})
	})

	// If no pronunciations found, add a default one
	if len(pronunciations) == 0 {
		pronunciations = append(pronunciations, models.Pronunciation{
			Word:     "n/a",
			Phonetic: "n/a",
			Audio:    "no_audio_available",
		})
	}

	return pronunciations
}
