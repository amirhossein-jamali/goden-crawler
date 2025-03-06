package repository

import (
	"errors"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/internal/db/elasticsearch"
	"github.com/amirhossein-jamali/goden-crawler/internal/db/mongodb"
	"github.com/amirhossein-jamali/goden-crawler/internal/db/postgres"
	"github.com/amirhossein-jamali/goden-crawler/internal/db/redis"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// WordRepository provides access to all database operations for words
type WordRepository struct {
	// Configuration options
	CacheTTL time.Duration
}

// NewWordRepository creates a new WordRepository
func NewWordRepository() *WordRepository {
	return &WordRepository{
		CacheTTL: 24 * time.Hour, // Default cache TTL
	}
}

// WithCacheTTL sets the cache TTL and returns the repository for chaining
func (r *WordRepository) WithCacheTTL(ttl time.Duration) *WordRepository {
	r.CacheTTL = ttl
	return r
}

// SaveWord saves a word to all databases
func (r *WordRepository) SaveWord(word *models.Word) error {
	var lastErr error

	// Save to MongoDB (primary storage)
	if err := mongodb.SaveWord(word); err != nil {
		logger.Error("Failed to save word to MongoDB", logger.F("word", word.Word), logger.F("error", err))
		lastErr = err
	}

	// Save to PostgreSQL (relational data)
	if err := postgres.SaveWord(word); err != nil {
		logger.Error("Failed to save word to PostgreSQL", logger.F("word", word.Word), logger.F("error", err))
		lastErr = err
	}

	// Cache in Redis
	if err := redis.CacheWord(word, r.CacheTTL); err != nil {
		logger.Error("Failed to cache word in Redis", logger.F("word", word.Word), logger.F("error", err))
		lastErr = err
	}

	// Index in Elasticsearch
	if err := elasticsearch.IndexWord(word); err != nil {
		logger.Error("Failed to index word in Elasticsearch", logger.F("word", word.Word), logger.F("error", err))
		lastErr = err
	}

	return lastErr
}

// GetWord retrieves a word from the fastest available source
func (r *WordRepository) GetWord(wordText string) (*models.Word, error) {
	var word *models.Word
	var err error

	// Try Redis first (fastest)
	word, err = redis.GetCachedWord(wordText)
	if err == nil && word != nil {
		logger.Info("Word retrieved from Redis cache", logger.F("word", wordText))
		return word, nil
	}

	// Try MongoDB next
	word, err = mongodb.GetWord(wordText)
	if err == nil && word != nil {
		logger.Info("Word retrieved from MongoDB", logger.F("word", wordText))
		// Cache the result in Redis for next time
		_ = redis.CacheWord(word, r.CacheTTL)
		return word, nil
	}

	// Try PostgreSQL as fallback
	word, err = postgres.GetWord(wordText)
	if err == nil && word != nil {
		logger.Info("Word retrieved from PostgreSQL", logger.F("word", wordText))
		// Cache the result in Redis for next time
		_ = redis.CacheWord(word, r.CacheTTL)
		return word, nil
	}

	// Not found in any database
	return nil, errors.New("word not found in any database")
}

// SearchWords searches for words in Elasticsearch
func (r *WordRepository) SearchWords(query string) ([]models.Word, error) {
	return elasticsearch.SearchWords(query)
}

// Close closes all database connections
func (r *WordRepository) Close() {
	_ = mongodb.Close()
	_ = postgres.Close()
	_ = redis.Close()
}
