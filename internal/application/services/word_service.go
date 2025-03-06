// File: internal/application/services/word_service.go

package services

import (
	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
	"github.com/amirhossein-jamali/goden-crawler/internal/repository"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// WordService provides operations for word data
type WordService struct {
	crawler    crawler.DudenCrawler
	repository *repository.WordRepository
}

// NewWordService creates a new WordService
func NewWordService(crawler crawler.DudenCrawler, repository *repository.WordRepository) *WordService {
	return &WordService{
		crawler:    crawler,
		repository: repository,
	}
}

// GetWordData retrieves word data from either the repository or by crawling
func (s *WordService) GetWordData(word string) (*models.Word, error) {
	logger.Info("Getting word data", logger.F("word", word))

	// First try to get the word from the repository
	wordData, err := s.repository.GetWord(word)
	if err == nil {
		logger.Info("Word found in repository", logger.F("word", word))
		return wordData, nil
	}

	// If not found in repository, crawl the data
	logger.Info("Word not found in repository, crawling", logger.F("word", word))
	wordData, err = s.crawler.FetchWordDataStructured(word)
	if err != nil {
		logger.Error("Failed to fetch word data", logger.F("word", word), logger.F("error", err))
		return nil, err
	}

	// Save the word data to the repository
	err = s.repository.SaveWord(wordData)
	if err != nil {
		logger.Error("Failed to save word to repository", logger.F("word", word), logger.F("error", err))
		// Continue even if saving fails
	}

	return wordData, nil
}

// GetWordSuggestions retrieves word suggestions
func (s *WordService) GetWordSuggestions(word string) ([]models.Synonym, error) {
	logger.Info("Getting word suggestions", logger.F("word", word))

	// First try to search in Elasticsearch
	words, err := s.repository.SearchWords(word)
	if err == nil && len(words) > 0 {
		synonyms := make([]models.Synonym, 0, len(words))
		for _, w := range words {
			synonyms = append(synonyms, models.Synonym{
				Text: w.Word,
			})
		}
		return synonyms, nil
	}

	// If not found or error, fall back to crawler
	return s.crawler.GetSuggestions(word)
}

// GetAvailableSections returns all available sections
func (s *WordService) GetAvailableSections() []string {
	return s.crawler.GetAvailableSections()
}

// Close closes the repository connections
func (s *WordService) Close() {
	s.repository.Close()
}
