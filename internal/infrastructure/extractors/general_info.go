// File: internal/infrastructure/extractors/general_info.go

package extractors

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
	"github.com/amirhossein-jamali/goden-crawler/pkg/utils"
)

// GeneralInfoExtractor extracts general information about a word
type GeneralInfoExtractor struct {
	*BaseExtractor
}

// init registers the extractor
func init() {
	RegisterExtractor("GeneralInfo", NewGeneralInfoExtractor)
}

// NewGeneralInfoExtractor creates a new GeneralInfoExtractor
func NewGeneralInfoExtractor(doc *goquery.Document) Extractor {
	return &GeneralInfoExtractor{
		BaseExtractor: NewBaseExtractor(doc, "GeneralInfo"),
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
		return []string{}
	}

	return utils.SplitAndTrim(wordTypeText, ",")
}

// extractFrequency extracts the frequency from the page
func (e *GeneralInfoExtractor) extractFrequency() string {
	frequencyElement := e.Doc.Find(".tuple__key:-soup-contains('Häufigkeit') + .tuple__val .shaft")
	if frequencyElement.Length() == 0 {
		return "unknown"
	}

	frequencyText := frequencyElement.Text()
	filledBars := strings.Count(frequencyText, "▒")

	switch filledBars {
	case 5:
		return "very_high"
	case 4:
		return "high"
	case 3:
		return "medium"
	case 2:
		return "low"
	case 1:
		return "very_low"
	default:
		return "unknown"
	}
}

// extractPronunciation extracts pronunciation information from the page
func (e *GeneralInfoExtractor) extractPronunciation() []models.Pronunciation {
	var pronunciations []models.Pronunciation
	e.Doc.Find(".pronunciation-guide").Each(func(i int, s *goquery.Selection) {
		word := s.Find(".pronunciation-guide__text").Text()
		phonetic := s.Find(".ipa").Text()
		audio := ""

		// Extract audio link if available
		audioElement := s.Find("a.pronunciation-guide__sound[data-duden-ref-type='audio']")
		if audioElement.Length() > 0 {
			audio, _ = audioElement.Attr("href")
		}

		pronunciations = append(pronunciations, models.Pronunciation{
			Word:     utils.CleanText(word),
			Phonetic: utils.CleanText(phonetic),
			Audio:    audio,
		})
	})

	if len(pronunciations) == 0 {
		return []models.Pronunciation{
			{
				Word:     "",
				Phonetic: "",
				Audio:    "",
			},
		}
	}

	return pronunciations
}
