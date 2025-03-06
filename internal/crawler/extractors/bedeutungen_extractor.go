// internal/crawler/extractors/bedeutungen_extractor.go
package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// BedeutungenExtractor extracts meanings and examples
type BedeutungenExtractor struct {
	*BaseExtractor
}

// NewBedeutungenExtractor creates a new BedeutungenExtractor
func NewBedeutungenExtractor(doc *goquery.Document) Extractor {
	return &BedeutungenExtractor{
		BaseExtractor: NewBaseExtractor(doc),
	}
}

// Extract extracts meanings and examples
func (e *BedeutungenExtractor) Extract() interface{} {
	return e.extractMeanings()
}

// extractMeanings extracts multiple meanings, grammar, examples, and images
func (e *BedeutungenExtractor) extractMeanings() []models.Meaning {
	var meanings []models.Meaning

	e.Doc.Find("#bedeutungen .enumeration__item").Each(func(i int, s *goquery.Selection) {
		// Extract parent meaning
		parentMeaning := e.extractSectionData(s)

		// Extract sub-meanings
		s.Find(".enumeration__sub-item").Each(func(j int, subSec *goquery.Selection) {
			subMeaning := e.extractSectionData(subSec)

			// If sub-meaning has the same text as parent, merge them
			if subMeaning.Text == parentMeaning.Text {
				e.mergeData(&parentMeaning, subMeaning)
			} else {
				parentMeaning.SubMeanings = append(parentMeaning.SubMeanings, subMeaning)
			}
		})

		meanings = append(meanings, parentMeaning)
	})

	// If no meanings found, add a default one
	if len(meanings) == 0 {
		meanings = append(meanings, models.Meaning{
			Text:         "n/a",
			Grammar:      "n/a",
			Examples:     []string{},
			Idioms:       []string{},
			Image:        "n/a",
			ImageCaption: "n/a",
			SubMeanings:  []models.Meaning{},
			TupleInfo:    map[string]string{},
		})
	}

	return meanings
}

// extractSectionData extracts data from a single meaning or sub-meaning section
func (e *BedeutungenExtractor) extractSectionData(section *goquery.Selection) models.Meaning {
	meaningText := section.Find(".enumeration__text").Text()
	meaningText = e.CleanText(meaningText)

	// Extract examples and idioms
	exampleData := e.extractExamples(section)

	// Extract image and caption
	var imageURL, imageCaption string
	section.Find(".depiction a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			imageURL = href
		}
	})

	section.Find(".depiction__caption").Each(func(i int, s *goquery.Selection) {
		imageCaption = e.CleanText(s.Text())
	})

	// Extract tuple info
	tupleInfo := e.parseTupleInfo(section)
	grammar, exists := tupleInfo["Grammatik"]
	if exists {
		delete(tupleInfo, "Grammatik")
	}

	if imageURL == "" {
		imageURL = "n/a"
	}
	if imageCaption == "" {
		imageCaption = "n/a"
	}
	if grammar == "" {
		grammar = "n/a"
	}

	return models.Meaning{
		Text:         meaningText,
		Grammar:      grammar,
		Examples:     exampleData["examples"],
		Idioms:       exampleData["idioms"],
		Image:        imageURL,
		ImageCaption: imageCaption,
		SubMeanings:  []models.Meaning{},
		TupleInfo:    tupleInfo,
	}
}

// parseTupleInfo parses all <dl class='tuple'> blocks in the current container
func (e *BedeutungenExtractor) parseTupleInfo(container *goquery.Selection) map[string]string {
	result := make(map[string]string)

	container.Find("dl.tuple").Each(func(i int, s *goquery.Selection) {
		s.Find("dt.tuple__key").Each(func(j int, dt *goquery.Selection) {
			dtText := e.CleanText(dt.Text())

			// Find the corresponding dd
			ddText := ""
			dt.NextAll().Each(func(k int, dd *goquery.Selection) {
				if dd.Is("dd.tuple__val") {
					ddText = e.CleanText(dd.Text())
					return
				}
			})

			if dtText != "" && ddText != "" {
				result[dtText] = ddText
			}
		})
	})

	return result
}

// extractExamples extracts both general examples and idioms
func (e *BedeutungenExtractor) extractExamples(section *goquery.Selection) map[string][]string {
	result := map[string][]string{
		"examples": {},
		"idioms":   {},
	}

	section.Find("dl.note").Each(func(i int, s *goquery.Selection) {
		title := s.Find("dt.note__title").Text()

		if title != "" && e.CleanText(title) == "Wendungen, Redensarten, Sprichwörter" {
			// Extract idioms
			s.Find("li").Each(func(j int, li *goquery.Selection) {
				idiom := e.CleanText(li.Text())
				if idiom != "" {
					result["idioms"] = append(result["idioms"], idiom)
				}
			})
		} else {
			// Extract general examples
			s.Find("li").Each(func(j int, li *goquery.Selection) {
				example := e.CleanText(li.Text())
				if example != "" {
					result["examples"] = append(result["examples"], example)
				}
			})
		}
	})

	return result
}

// mergeData merges data from sub-meaning into the parent meaning
func (e *BedeutungenExtractor) mergeData(parent *models.Meaning, sub models.Meaning) {
	// Merge examples
	for _, example := range sub.Examples {
		// Check if example already exists
		exists := false
		for _, parentExample := range parent.Examples {
			if parentExample == example {
				exists = true
				break
			}
		}
		if !exists {
			parent.Examples = append(parent.Examples, example)
		}
	}

	// Merge idioms
	for _, idiom := range sub.Idioms {
		// Check if idiom already exists
		exists := false
		for _, parentIdiom := range parent.Idioms {
			if parentIdiom == idiom {
				exists = true
				break
			}
		}
		if !exists {
			parent.Idioms = append(parent.Idioms, idiom)
		}
	}

	// If parent's grammar is n/a but sub has grammar, update parent's
	if (parent.Grammar == "" || parent.Grammar == "n/a") && sub.Grammar != "n/a" {
		parent.Grammar = sub.Grammar
	}

	// If parent image is n/a but sub has an image
	if (parent.Image == "" || parent.Image == "n/a") && sub.Image != "n/a" {
		parent.Image = sub.Image
		parent.ImageCaption = sub.ImageCaption
	}

	// Merge tuple_info
	for k, v := range sub.TupleInfo {
		if _, exists := parent.TupleInfo[k]; !exists || parent.TupleInfo[k] == "" {
			parent.TupleInfo[k] = v
		}
	}
}
