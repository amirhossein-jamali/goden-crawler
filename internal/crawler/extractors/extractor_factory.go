// File: internal/crawler/extractors/extractor_factory.go

package extractors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// ExtractorFactory creates and manages extractors
type ExtractorFactory struct {
	extractors map[string]func(*goquery.Document) Extractor
}

// NewExtractorFactory creates a new ExtractorFactory
func NewExtractorFactory() *ExtractorFactory {
	factory := &ExtractorFactory{
		extractors: make(map[string]func(*goquery.Document) Extractor),
	}

	// Register all extractors
	factory.registerExtractors()

	return factory
}

// registerExtractors registers all available extractors
func (f *ExtractorFactory) registerExtractors() {
	// Register general info extractor
	f.extractors["general_info"] = func(doc *goquery.Document) Extractor {
		return NewGeneralInfoExtractor(doc)
	}

	// Register meanings extractor
	f.extractors["bedeutungen"] = func(doc *goquery.Document) Extractor {
		return NewBedeutungenExtractor(doc)
	}

	// Register grammar extractor
	f.extractors["grammatik"] = func(doc *goquery.Document) Extractor {
		return NewGrammatikExtractor(doc)
	}

	// Register spelling extractor
	f.extractors["rechtschreibung"] = func(doc *goquery.Document) Extractor {
		return NewRechtschreibungExtractor(doc)
	}

	// Register synonyms extractor
	f.extractors["synonyme"] = func(doc *goquery.Document) Extractor {
		return NewSynonymeExtractor(doc)
	}

	// Register origin extractor
	f.extractors["herkunft"] = func(doc *goquery.Document) Extractor {
		return NewHerkunftExtractor(doc)
	}

	// Register fun facts extractor
	f.extractors["wussten_sie_schon"] = func(doc *goquery.Document) Extractor {
		return NewWusstenSieSchonExtractor(doc)
	}
}

// CreateExtractor creates an extractor for the given section
func (f *ExtractorFactory) CreateExtractor(section string, doc *goquery.Document) (Extractor, error) {
	constructor, exists := f.extractors[section]
	if !exists {
		return nil, fmt.Errorf("unknown extractor for section: %s", section)
	}

	return constructor(doc), nil
}

// GetAvailableSections returns all available section names
func (f *ExtractorFactory) GetAvailableSections() []string {
	sections := make([]string, 0, len(f.extractors))
	for section := range f.extractors {
		sections = append(sections, section)
	}
	return sections
}
