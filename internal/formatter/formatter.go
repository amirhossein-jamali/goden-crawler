// ./internal/formatter/formatter.go
package formatter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
		return formatAsText(wordData), nil
	default:
		return "", errors.New("invalid format selected")
	}
}

// formatAsText formats the word data as human-readable text
func formatAsText(wordData *models.Word) string {
	var sb strings.Builder

	// General information
	sb.WriteString(fmt.Sprintf("Word: %s\n", wordData.Word))
	if wordData.Article != "" {
		sb.WriteString(fmt.Sprintf("Article: %s\n", wordData.Article))
	}
	if len(wordData.WordType) > 0 {
		sb.WriteString(fmt.Sprintf("Word Type: %s\n", strings.Join(wordData.WordType, ", ")))
	}
	if wordData.Frequency != "" {
		sb.WriteString(fmt.Sprintf("Frequency: %s\n", wordData.Frequency))
	}
	if wordData.Grammar != "" {
		sb.WriteString(fmt.Sprintf("Grammar: %s\n", wordData.Grammar))
	}

	// Pronunciation
	if len(wordData.Pronunciation) > 0 {
		sb.WriteString("\nPronunciation:\n")
		for _, p := range wordData.Pronunciation {
			sb.WriteString(fmt.Sprintf("  Word: %s\n", p.Word))
			sb.WriteString(fmt.Sprintf("  Phonetic: %s\n", p.Phonetic))
			if p.Audio != "" && p.Audio != "no_audio_available" {
				sb.WriteString(fmt.Sprintf("  Audio: %s\n", p.Audio))
			}
		}
	}

	// Meanings
	if len(wordData.Meanings) > 0 {
		sb.WriteString("\nMeanings:\n")
		for i, meaning := range wordData.Meanings {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, meaning.Text))

			if meaning.Grammar != "" && meaning.Grammar != "n/a" {
				sb.WriteString(fmt.Sprintf("   Grammar: %s\n", meaning.Grammar))
			}

			if len(meaning.Examples) > 0 {
				sb.WriteString("   Examples:\n")
				for _, example := range meaning.Examples {
					sb.WriteString(fmt.Sprintf("     - %s\n", example))
				}
			}

			if len(meaning.Idioms) > 0 {
				sb.WriteString("   Idioms:\n")
				for _, idiom := range meaning.Idioms {
					sb.WriteString(fmt.Sprintf("     - %s\n", idiom))
				}
			}

			if len(meaning.SubMeanings) > 0 {
				sb.WriteString("   Sub-meanings:\n")
				for j, subMeaning := range meaning.SubMeanings {
					sb.WriteString(fmt.Sprintf("     %d.%d. %s\n", i+1, j+1, subMeaning.Text))
					if len(subMeaning.Examples) > 0 {
						for _, example := range subMeaning.Examples {
							sb.WriteString(fmt.Sprintf("       - %s\n", example))
						}
					}
				}
			}
		}
	}

	// Synonyms
	if len(wordData.Synonyms) > 0 {
		sb.WriteString("\nSynonyms:\n")
		for _, synonym := range wordData.Synonyms {
			sb.WriteString(fmt.Sprintf("  - %s\n", synonym.Text))
		}
	}

	// Origin
	if len(wordData.Origin) > 0 {
		sb.WriteString("\nOrigin:\n")
		var originText string
		for _, origin := range wordData.Origin {
			originText += origin.Word + " "
		}
		sb.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(originText)))
	}

	// Fun Facts
	if len(wordData.FunFacts) > 0 {
		sb.WriteString("\nDid you know?\n")
		for _, fact := range wordData.FunFacts {
			sb.WriteString(fmt.Sprintf("  - %s\n", fact))
		}
	}

	return sb.String()
}
