// ./internal/formatter/formatter.go
package formatter

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// FormatOutput formats the output based on user selection
func FormatOutput(wordData *models.Word, format string) (string, error) {
	switch format {
	case "json":
		jsonData, err := json.MarshalIndent(wordData, "", "  ")
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	case "text":
		return fmt.Sprintf(
			"Word: %s\nMeanings: %v\nSynonyms: %v\nGrammar: %s",
			wordData.Word, wordData.Meanings, wordData.Synonyms, wordData.Grammar,
		), nil
	default:
		return "", errors.New("invalid format selected")
	}
}
