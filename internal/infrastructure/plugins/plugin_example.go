// File: internal/infrastructure/plugins/plugin_example.go

package plugins

import (
	"encoding/xml"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/internal/domain/interfaces"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// This file serves as an example of how to create a plugin for the Goden Crawler.
// To create a real plugin, you would compile this as a separate Go plugin (.so file).

// XMLFormatter is an example formatter that outputs word data as XML
type XMLFormatter struct{}

// FormatOutput formats the word data as XML
func (f *XMLFormatter) FormatOutput(wordData *models.Word, format string) (string, error) {
	if format != "xml" {
		return "", fmt.Errorf("XMLFormatter only supports 'xml' format")
	}

	xmlData, err := xml.MarshalIndent(wordData, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(xmlData), nil
}

// CustomExtractor is an example extractor that extracts custom data
type CustomExtractor struct {
	doc  *goquery.Document
	name string
}

// NewCustomExtractor creates a new CustomExtractor
func NewCustomExtractor(doc *goquery.Document) interfaces.Extractor {
	return &CustomExtractor{
		doc:  doc,
		name: "custom",
	}
}

// Extract extracts custom data from the document
func (e *CustomExtractor) Extract() (interface{}, error) {
	// This is just an example implementation
	// In a real plugin, you would extract actual data from the document
	return "Custom data extracted", nil
}

// GetName returns the name of the extractor
func (e *CustomExtractor) GetName() string {
	return e.name
}

// Init is the entry point for the plugin
// This function is called when the plugin is loaded
func Init(pm *PluginManager) error {
	// Register the XML formatter
	pm.RegisterFormatter("xml", &XMLFormatter{})

	// Register the custom extractor factory
	pm.RegisterExtractor("custom", &CustomExtractor{
		doc:  nil, // This will be set when the extractor is used
		name: "custom",
	})

	return nil
}

/*
To compile this as a plugin:

1. Create a separate Go file with the following content:

```go
package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/internal/domain/interfaces"
	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/plugins"
)

// XMLFormatter implementation...
// CustomExtractor implementation...

// Init is exported as a symbol for the plugin system
func Init(pm *plugins.PluginManager) error {
	// Register formatters and extractors
	return nil
}
```

2. Compile the plugin:

```bash
go build -buildmode=plugin -o xml_formatter.so xml_formatter.go
```

3. Place the compiled plugin in the plugins directory.
*/
