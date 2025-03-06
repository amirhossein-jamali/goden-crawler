// File: internal/infrastructure/extractors/factory.go

package extractors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	customErrors "github.com/amirhossein-jamali/goden-crawler/pkg/errors"
	"github.com/amirhossein-jamali/goden-crawler/pkg/utils"
)

// ExtractorFactory creates and manages extractors
type ExtractorFactory struct {
	// No need to store extractors here, they're in the registry
}

// NewExtractorFactory creates a new ExtractorFactory
func NewExtractorFactory() *ExtractorFactory {
	// Register all extractors
	registerAllExtractors()

	return &ExtractorFactory{}
}

// registerAllExtractors ensures all extractors are registered
// This is called automatically when creating a new factory
func registerAllExtractors() {
	// The registration is done in the init() functions of each extractor file
	// This function exists to document the registration process
}

// CreateExtractor creates an extractor for the given section
func (f *ExtractorFactory) CreateExtractor(section string, doc *goquery.Document) (Extractor, error) {
	// Convert section name to extractor name (e.g., "general_info" -> "GeneralInfo")
	extractorName := utils.ToTitleCase(section)

	constructor, exists := GetExtractor(extractorName)
	if !exists {
		return nil, customErrors.NewExtractorError(
			section,
			fmt.Sprintf("unknown extractor: %s", extractorName),
			nil,
		)
	}

	return constructor(doc), nil
}

// GetAvailableSections returns all available section names
func (f *ExtractorFactory) GetAvailableSections() []string {
	extractors := GetAllExtractors()
	sections := make([]string, 0, len(extractors))

	for name := range extractors {
		// Convert extractor name to section name (e.g., "GeneralInfo" -> "general_info")
		sectionName := utils.ToSnakeCase(name)
		sections = append(sections, sectionName)
	}

	return sections
}
