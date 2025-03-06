// File: internal/application/services/word_service.go

package services

import (
	"fmt"

	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/events"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// WordService provides methods for retrieving and processing word data
type WordService struct {
	crawler crawler.DudenCrawler
}

// NewWordService creates a new WordService with the given crawler
func NewWordService(crawler crawler.DudenCrawler) *WordService {
	return &WordService{
		crawler: crawler,
	}
}

// GetWordData fetches structured data for a word
func (s *WordService) GetWordData(word string) (*models.Word, error) {
	logger.Info("Fetching word data", logger.F("word", word))

	// Notify that word fetch is starting
	events.NotifyObservers(events.Event{
		Type: events.WordFetchStarted,
		Payload: map[string]string{
			"word": word,
		},
	})

	// Fetch word data using the crawler
	wordData, err := s.crawler.FetchWordDataStructured(word)
	if err != nil {
		logger.Error("Failed to fetch word data",
			logger.F("word", word),
			logger.F("error", err.Error()))

		// Notify that word fetch failed
		events.NotifyObservers(events.Event{
			Type: events.WordFetchFailed,
			Payload: map[string]interface{}{
				"word":  word,
				"error": err.Error(),
			},
		})

		return nil, fmt.Errorf("failed to fetch word data: %w", err)
	}

	// Notify that word fetch completed
	events.NotifyObservers(events.Event{
		Type: events.WordFetchCompleted,
		Payload: map[string]interface{}{
			"word": word,
			"data": wordData,
		},
	})

	logger.Info("Successfully fetched word data", logger.F("word", word))
	return wordData, nil
}

// GetWordSuggestions fetches suggestions for a word
func (s *WordService) GetWordSuggestions(word string) ([]models.Synonym, error) {
	logger.Info("Fetching word suggestions", logger.F("word", word))

	suggestions, err := s.crawler.GetSuggestions(word)
	if err != nil {
		logger.Error("Failed to fetch word suggestions",
			logger.F("word", word),
			logger.F("error", err.Error()))
		return nil, fmt.Errorf("failed to fetch word suggestions: %w", err)
	}

	logger.Info("Found suggestions",
		logger.F("word", word),
		logger.F("count", len(suggestions)))
	return suggestions, nil
}

// GetAvailableSections returns all available sections that can be extracted
func (s *WordService) GetAvailableSections() []string {
	return s.crawler.GetAvailableSections()
}
