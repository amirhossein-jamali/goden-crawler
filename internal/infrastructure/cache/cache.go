// File: internal/infrastructure/cache/cache.go

package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// Cache provides caching functionality for word data
type Cache struct {
	cachePath    string
	cacheTTL     time.Duration
	memoryCache  map[string]*cacheEntry
	mu           sync.RWMutex
	enableMemory bool
	enableDisk   bool
}

// cacheEntry represents a cached item with metadata
type cacheEntry struct {
	Data      *models.Word `json:"data"`
	Timestamp time.Time    `json:"timestamp"`
}

// NewCache creates a new cache instance
func NewCache(cachePath string, cacheTTL time.Duration, enableMemory, enableDisk bool) *Cache {
	// Create cache directory if it doesn't exist
	if enableDisk && cachePath != "" {
		if err := os.MkdirAll(cachePath, 0755); err != nil {
			logger.Warn("Failed to create cache directory",
				logger.F("path", cachePath),
				logger.F("error", err))
		}
	}

	return &Cache{
		cachePath:    cachePath,
		cacheTTL:     cacheTTL,
		memoryCache:  make(map[string]*cacheEntry),
		enableMemory: enableMemory,
		enableDisk:   enableDisk,
	}
}

// Get retrieves a word from the cache
func (c *Cache) Get(word string) (*models.Word, bool) {
	// Try memory cache first
	if c.enableMemory {
		c.mu.RLock()
		entry, found := c.memoryCache[word]
		c.mu.RUnlock()

		if found && time.Since(entry.Timestamp) < c.cacheTTL {
			logger.Debug("Cache hit (memory)", logger.F("word", word))
			return entry.Data, true
		}
	}

	// Try disk cache if memory cache failed
	if c.enableDisk {
		data, found := c.getFromDisk(word)
		if found {
			// Update memory cache
			if c.enableMemory {
				c.mu.Lock()
				c.memoryCache[word] = &cacheEntry{
					Data:      data,
					Timestamp: time.Now(),
				}
				c.mu.Unlock()
			}
			return data, true
		}
	}

	return nil, false
}

// Set stores a word in the cache
func (c *Cache) Set(word string, data *models.Word) {
	// Update memory cache
	if c.enableMemory {
		c.mu.Lock()
		c.memoryCache[word] = &cacheEntry{
			Data:      data,
			Timestamp: time.Now(),
		}
		c.mu.Unlock()
	}

	// Update disk cache
	if c.enableDisk {
		c.saveToDisk(word, data)
	}
}

// Clear clears the cache
func (c *Cache) Clear() error {
	// Clear memory cache
	if c.enableMemory {
		c.mu.Lock()
		c.memoryCache = make(map[string]*cacheEntry)
		c.mu.Unlock()
	}

	// Clear disk cache
	if c.enableDisk && c.cachePath != "" {
		files, err := filepath.Glob(filepath.Join(c.cachePath, "*.json"))
		if err != nil {
			return fmt.Errorf("failed to list cache files: %w", err)
		}

		for _, file := range files {
			if err := os.Remove(file); err != nil {
				logger.Warn("Failed to remove cache file",
					logger.F("file", file),
					logger.F("error", err))
			}
		}
	}

	return nil
}

// getFromDisk retrieves a word from the disk cache
func (c *Cache) getFromDisk(word string) (*models.Word, bool) {
	if c.cachePath == "" {
		return nil, false
	}

	filePath := filepath.Join(c.cachePath, fmt.Sprintf("%s.json", word))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		logger.Warn("Failed to unmarshal cache entry",
			logger.F("word", word),
			logger.F("error", err))
		return nil, false
	}

	// Check if the entry is expired
	if time.Since(entry.Timestamp) > c.cacheTTL {
		logger.Debug("Cache entry expired", logger.F("word", word))
		return nil, false
	}

	logger.Debug("Cache hit (disk)", logger.F("word", word))
	return entry.Data, true
}

// saveToDisk stores a word in the disk cache
func (c *Cache) saveToDisk(word string, data *models.Word) {
	if c.cachePath == "" {
		return
	}

	entry := cacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		logger.Warn("Failed to marshal cache entry",
			logger.F("word", word),
			logger.F("error", err))
		return
	}

	filePath := filepath.Join(c.cachePath, fmt.Sprintf("%s.json", word))
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		logger.Warn("Failed to write cache file",
			logger.F("word", word),
			logger.F("path", filePath),
			logger.F("error", err))
	}
}
