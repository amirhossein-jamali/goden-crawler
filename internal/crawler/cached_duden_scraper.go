// File: internal/crawler/cached_duden_scraper.go

package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/cache"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// CachedDudenScraper wraps a DudenScraper with caching functionality
type CachedDudenScraper struct {
	scraper *DudenScraper
	cache   *cache.Cache
}

// NewCachedDudenScraper creates a new CachedDudenScraper
func NewCachedDudenScraper(scraper *DudenScraper, cache *cache.Cache) *CachedDudenScraper {
	return &CachedDudenScraper{
		scraper: scraper,
		cache:   cache,
	}
}

// makeRequest delegates to the underlying scraper's makeRequest method
func (s *CachedDudenScraper) makeRequest(url string) (*goquery.Document, error) {
	// No caching for HTML documents
	return s.scraper.makeRequest(url)
}

// fetchWordDoc delegates to the underlying scraper's fetchWordDoc method
func (s *CachedDudenScraper) fetchWordDoc(word string) (*goquery.Document, error) {
	// No caching for HTML documents
	return s.scraper.fetchWordDoc(word)
}

// FetchWordDataStructured fetches structured data for a word
func (s *CachedDudenScraper) FetchWordDataStructured(word string) (*models.Word, error) {
	// Try to get from cache first
	cachedData, found := s.cache.Get(word)
	if found {
		logger.Info("Using cached data for word", logger.F("word", word))
		return cachedData, nil
	}

	// If not in cache, fetch from source
	logger.Info("Fetching word data from source", logger.F("word", word))
	data, err := s.scraper.FetchWordDataStructured(word)
	if err != nil {
		return nil, err
	}

	// Store in cache
	s.cache.Set(word, data)
	return data, nil
}

// GetSuggestions returns a list of suggested words for a given input
func (s *CachedDudenScraper) GetSuggestions(word string) ([]models.Synonym, error) {
	// No caching for suggestions as they might change
	return s.scraper.GetSuggestions(word)
}

// GetAvailableSections returns a list of all available sections that can be extracted
func (s *CachedDudenScraper) GetAvailableSections() []string {
	return s.scraper.GetAvailableSections()
}
