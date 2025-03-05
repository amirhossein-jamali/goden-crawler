// ./pkg/word.go
package models

// Word represents the structure for storing word data
type Word struct {
	Word     string   `json:"word"`
	Meanings []string `json:"meanings"`
	Synonyms []string `json:"synonyms"`
	Grammar  string   `json:"grammar"`
}
