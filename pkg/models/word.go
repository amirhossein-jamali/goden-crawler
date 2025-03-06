// ./pkg/models/word.go
package models

// Word represents a German word with its linguistic information
type Word struct {
	Word          string          `json:"word"`
	Article       string          `json:"article,omitempty"`
	WordType      []string        `json:"word_type,omitempty"`
	Frequency     string          `json:"frequency,omitempty"`
	Grammar       string          `json:"grammar,omitempty"`
	Meanings      []Meaning       `json:"meanings,omitempty"`
	Synonyms      []Synonym       `json:"synonyms,omitempty"`
	Pronunciation []Pronunciation `json:"pronunciation,omitempty"`
	Spelling      Spelling        `json:"spelling,omitempty"`
	Origin        []Origin        `json:"origin,omitempty"`
	FunFacts      []string        `json:"fun_facts,omitempty"`
}

// Meaning represents a single meaning of a word
type Meaning struct {
	Text         string            `json:"text"`
	Grammar      string            `json:"grammar,omitempty"`
	Examples     []string          `json:"examples,omitempty"`
	Idioms       []string          `json:"idioms,omitempty"`
	Image        string            `json:"image,omitempty"`
	ImageCaption string            `json:"image_caption,omitempty"`
	SubMeanings  []Meaning         `json:"sub_meanings,omitempty"`
	TupleInfo    map[string]string `json:"tuple_info,omitempty"`
}

// Synonym represents a synonym with optional link
type Synonym struct {
	Text string `json:"text"`
	Link string `json:"link,omitempty"`
}

// Pronunciation represents pronunciation information
type Pronunciation struct {
	Word     string `json:"word"`
	Phonetic string `json:"phonetic"`
	Audio    string `json:"audio,omitempty"`
}

// Spelling represents spelling information
type Spelling struct {
	SyllabicDivision string   `json:"syllabic_division,omitempty"`
	Examples         []string `json:"examples,omitempty"`
	Rules            []Rule   `json:"rules,omitempty"`
}

// Rule represents a grammatical rule with link
type Rule struct {
	Text string `json:"text"`
	Link string `json:"link"`
}

// Origin represents word origin information
type Origin struct {
	Word string `json:"word"`
	Link string `json:"link,omitempty"`
}
