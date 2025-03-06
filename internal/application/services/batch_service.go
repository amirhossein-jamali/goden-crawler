// File: internal/application/services/batch_service.go

package services

import (
	"context"
	"sync"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// BatchResult represents the result of processing a single word
type BatchResult struct {
	Word  string
	Data  *models.Word
	Error error
}

// BatchService provides methods for batch processing of words
type BatchService struct {
	wordService *WordService
	workers     int
	timeout     time.Duration
}

// NewBatchService creates a new BatchService
func NewBatchService(crawler crawler.DudenCrawler, workers int, timeout time.Duration) *BatchService {
	return &BatchService{
		wordService: NewWordService(crawler),
		workers:     workers,
		timeout:     timeout,
	}
}

// ProcessWords processes multiple words concurrently
func (s *BatchService) ProcessWords(ctx context.Context, words []string) []BatchResult {
	results := make([]BatchResult, 0, len(words))
	resultChan := make(chan BatchResult, len(words))

	// Create a worker pool
	var wg sync.WaitGroup
	wordChan := make(chan string, len(words))

	// Start workers
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.worker(ctx, &wg, wordChan, resultChan)
	}

	// Send words to workers
	for _, word := range words {
		wordChan <- word
	}
	close(wordChan)

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Collect results
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// worker processes words from the word channel
func (s *BatchService) worker(ctx context.Context, wg *sync.WaitGroup, wordChan <-chan string, resultChan chan<- BatchResult) {
	defer wg.Done()

	for {
		select {
		case word, ok := <-wordChan:
			if !ok {
				return
			}

			// Create a context with timeout
			workerCtx, cancel := context.WithTimeout(ctx, s.timeout)

			// Process the word
			result := s.processWord(workerCtx, word)
			resultChan <- result

			cancel()

		case <-ctx.Done():
			logger.Info("Worker cancelled", logger.F("reason", ctx.Err()))
			return
		}
	}
}

// processWord processes a single word
func (s *BatchService) processWord(ctx context.Context, word string) BatchResult {
	logger.Info("Processing word", logger.F("word", word))

	// Create a channel to receive the result
	resultChan := make(chan BatchResult, 1)

	// Process the word in a goroutine
	go func() {
		data, err := s.wordService.GetWordData(word)
		resultChan <- BatchResult{
			Word:  word,
			Data:  data,
			Error: err,
		}
	}()

	// Wait for the result or timeout
	select {
	case result := <-resultChan:
		return result
	case <-ctx.Done():
		logger.Warn("Word processing timed out", logger.F("word", word))
		return BatchResult{
			Word:  word,
			Data:  nil,
			Error: ctx.Err(),
		}
	}
}
