package models

import (
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// WordBuilder is a builder for Word objects
type WordBuilder struct {
	word          string
	article       string
	wordType      []string
	frequency     string
	grammar       string
	meanings      []models.Meaning
	synonyms      []models.Synonym
	pronunciation []models.Pronunciation
	spelling      models.Spelling
	origin        []models.Origin
	funFacts      []string
}

// NewWordBuilder creates a new WordBuilder
func NewWordBuilder() *WordBuilder {
	return &WordBuilder{
		wordType:      []string{},
		meanings:      []models.Meaning{},
		synonyms:      []models.Synonym{},
		pronunciation: []models.Pronunciation{},
		origin:        []models.Origin{},
		funFacts:      []string{},
	}
}

// WithWord sets the word
func (b *WordBuilder) WithWord(word string) *WordBuilder {
	b.word = word
	return b
}

// WithArticle sets the article
func (b *WordBuilder) WithArticle(article string) *WordBuilder {
	b.article = article
	return b
}

// WithWordType sets the word type
func (b *WordBuilder) WithWordType(wordType []string) *WordBuilder {
	b.wordType = wordType
	return b
}

// AddWordType adds a word type
func (b *WordBuilder) AddWordType(wordType string) *WordBuilder {
	b.wordType = append(b.wordType, wordType)
	return b
}

// WithFrequency sets the frequency
func (b *WordBuilder) WithFrequency(frequency string) *WordBuilder {
	b.frequency = frequency
	return b
}

// WithGrammar sets the grammar
func (b *WordBuilder) WithGrammar(grammar string) *WordBuilder {
	b.grammar = grammar
	return b
}

// WithMeanings sets the meanings
func (b *WordBuilder) WithMeanings(meanings []models.Meaning) *WordBuilder {
	b.meanings = meanings
	return b
}

// AddMeaning adds a meaning
func (b *WordBuilder) AddMeaning(meaning models.Meaning) *WordBuilder {
	b.meanings = append(b.meanings, meaning)
	return b
}

// WithSynonyms sets the synonyms
func (b *WordBuilder) WithSynonyms(synonyms []models.Synonym) *WordBuilder {
	b.synonyms = synonyms
	return b
}

// AddSynonym adds a synonym
func (b *WordBuilder) AddSynonym(synonym models.Synonym) *WordBuilder {
	b.synonyms = append(b.synonyms, synonym)
	return b
}

// WithPronunciation sets the pronunciation
func (b *WordBuilder) WithPronunciation(pronunciation []models.Pronunciation) *WordBuilder {
	b.pronunciation = pronunciation
	return b
}

// AddPronunciation adds a pronunciation
func (b *WordBuilder) AddPronunciation(pronunciation models.Pronunciation) *WordBuilder {
	b.pronunciation = append(b.pronunciation, pronunciation)
	return b
}

// WithSpelling sets the spelling
func (b *WordBuilder) WithSpelling(spelling models.Spelling) *WordBuilder {
	b.spelling = spelling
	return b
}

// WithOrigin sets the origin
func (b *WordBuilder) WithOrigin(origin []models.Origin) *WordBuilder {
	b.origin = origin
	return b
}

// AddOrigin adds an origin
func (b *WordBuilder) AddOrigin(origin models.Origin) *WordBuilder {
	b.origin = append(b.origin, origin)
	return b
}

// WithFunFacts sets the fun facts
func (b *WordBuilder) WithFunFacts(funFacts []string) *WordBuilder {
	b.funFacts = funFacts
	return b
}

// AddFunFact adds a fun fact
func (b *WordBuilder) AddFunFact(funFact string) *WordBuilder {
	b.funFacts = append(b.funFacts, funFact)
	return b
}

// Build builds the Word object
func (b *WordBuilder) Build() *models.Word {
	return &models.Word{
		Word:          b.word,
		Article:       b.article,
		WordType:      b.wordType,
		Frequency:     b.frequency,
		Grammar:       b.grammar,
		Meanings:      b.meanings,
		Synonyms:      b.synonyms,
		Pronunciation: b.pronunciation,
		Spelling:      b.spelling,
		Origin:        b.origin,
		FunFacts:      b.funFacts,
	}
}
