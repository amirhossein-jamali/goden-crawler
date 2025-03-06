// File: internal/application/services/batch_service.go

package services

import (
	"context"
	"sync"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/internal/repository"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// BatchResult represents the result of a batch operation
type BatchResult struct {
	Word  string
	Data  *models.Word
	Error error
}

// BatchService provides batch processing operations
type BatchService struct {
	wordService *WordService
	workers     int
	timeout     time.Duration
	repository  *repository.WordRepository
}

// NewBatchService creates a new BatchService
func NewBatchService(wordService *WordService, repository *repository.WordRepository, workers int, timeout time.Duration) *BatchService {
	return &BatchService{
		wordService: wordService,
		repository:  repository,
		workers:     workers,
		timeout:     timeout,
	}
}

// ProcessWords processes multiple words in parallel
func (s *BatchService) ProcessWords(ctx context.Context, words []string) []BatchResult {
	logger.Info("Starting batch processing", logger.F("word_count", len(words)), logger.F("workers", s.workers))

	// Create channels for words and results
	wordChan := make(chan string, len(words))
	resultChan := make(chan BatchResult, len(words))

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.worker(ctx, &wg, wordChan, resultChan)
	}

	// Send words to the channel
	for _, word := range words {
		wordChan <- word
	}
	close(wordChan)

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Collect results
	var results []BatchResult
	for result := range resultChan {
		results = append(results, result)
	}

	logger.Info("Batch processing completed", logger.F("word_count", len(words)), logger.F("result_count", len(results)))
	return results
}

// worker processes words from the channel
func (s *BatchService) worker(ctx context.Context, wg *sync.WaitGroup, wordChan <-chan string, resultChan chan<- BatchResult) {
	defer wg.Done()

	for word := range wordChan {
		// Create a context with timeout
		tctx, cancel := context.WithTimeout(ctx, s.timeout)

		// Process the word
		result := s.processWord(tctx, word)
		resultChan <- result

		// Save the result to the repository if successful
		if result.Error == nil && result.Data != nil {
			err := s.repository.SaveWord(result.Data)
			if err != nil {
				logger.Error("Failed to save word to repository",
					logger.F("word", word),
					logger.F("error", err))
			}
		}

		cancel()
	}
}

// processWord processes a single word
func (s *BatchService) processWord(ctx context.Context, word string) BatchResult {
	logger.Info("Processing word", logger.F("word", word))

	// Create a channel to receive the result
	resultChan := make(chan BatchResult, 1)

	// Process the word in a goroutine
	go func() {
		// First check if the word already exists in the repository
		data, err := s.repository.GetWord(word)
		if err == nil && data != nil {
			logger.Info("Word found in repository", logger.F("word", word))
			resultChan <- BatchResult{
				Word:  word,
				Data:  data,
				Error: nil,
			}
			return
		}

		// If not found, get it from the word service
		data, err = s.wordService.GetWordData(word)
		resultChan <- BatchResult{
			Word:  word,
			Data:  data,
			Error: err,
		}
	}()

	// Wait for the result or timeout
	select {
	case <-ctx.Done():
		return BatchResult{
			Word:  word,
			Error: ctx.Err(),
		}
	case result := <-resultChan:
		return result
	}
}
