// File: internal/infrastructure/container/container.go

package container

import (
	"fmt"
	"sync"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/internal/application/services"
	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
	"github.com/amirhossein-jamali/goden-crawler/internal/infrastructure/cache"
	"github.com/amirhossein-jamali/goden-crawler/internal/repository"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
	"github.com/amirhossein-jamali/goden-crawler/pkg/utils"
)

// Container is a simple dependency injection container
type Container struct {
	services map[string]interface{}
	mutex    sync.RWMutex
}

// NewContainer creates a new container
func NewContainer() *Container {
	c := &Container{
		services: make(map[string]interface{}),
	}
	registerDefaultServices(c)
	return c
}

// Register registers a service with the container
func (c *Container) Register(name string, service interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.services[name] = service
	logger.Debug("Registered service", logger.F("name", name))
}

// Get retrieves a service from the container
func (c *Container) Get(name string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service not found: %s", name)
	}

	return service, nil
}

// MustGet retrieves a service or panics if it doesn't exist
func (c *Container) MustGet(name string) interface{} {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// GetWordService returns the WordService from the container
func (c *Container) GetWordService() *services.WordService {
	service, _ := c.Get("wordService")
	if service == nil {
		dudenScraper := c.GetCachedDudenScraper()
		wordRepo := c.GetWordRepository()
		wordService := services.NewWordService(dudenScraper, wordRepo)
		c.Register("wordService", wordService)
		return wordService
	}
	return service.(*services.WordService)
}

// GetDudenScraper returns the DudenScraper
func (c *Container) GetDudenScraper() *crawler.DudenScraper {
	service, _ := c.Get("dudenScraper")
	if service == nil {
		dudenScraper := crawler.NewDudenScraper()
		c.Register("dudenScraper", dudenScraper)
		return dudenScraper
	}
	return service.(*crawler.DudenScraper)
}

// GetCache returns the Cache
func (c *Container) GetCache() *cache.Cache {
	service, _ := c.Get("cache")
	if service == nil {
		// Create a new cache with default settings
		// Use memory cache only for now
		cacheInstance := cache.NewCache("", 24*time.Hour, true, false)
		c.Register("cache", cacheInstance)
		return cacheInstance
	}
	return service.(*cache.Cache)
}

// GetCachedDudenScraper returns the CachedDudenScraper
func (c *Container) GetCachedDudenScraper() *crawler.CachedDudenScraper {
	service, _ := c.Get("cachedDudenScraper")
	if service == nil {
		dudenScraper := c.GetDudenScraper()
		cacheInstance := c.GetCache()
		cachedDudenScraper := crawler.NewCachedDudenScraper(dudenScraper, cacheInstance)
		c.Register("cachedDudenScraper", cachedDudenScraper)
		return cachedDudenScraper
	}
	return service.(*crawler.CachedDudenScraper)
}

// GetWordRepository returns the WordRepository
func (c *Container) GetWordRepository() *repository.WordRepository {
	service, _ := c.Get("wordRepository")
	if service == nil {
		wordRepo := repository.NewWordRepository()
		c.Register("wordRepository", wordRepo)
		return wordRepo
	}
	return service.(*repository.WordRepository)
}

// GetBatchService returns the BatchService
func (c *Container) GetBatchService() *services.BatchService {
	service, _ := c.Get("batchService")
	if service == nil {
		wordService := c.GetWordService()
		wordRepo := c.GetWordRepository()
		batchService := services.NewBatchService(wordService, wordRepo, 5, 30*time.Second)
		c.Register("batchService", batchService)
		return batchService
	}
	return service.(*services.BatchService)
}

// Singleton instance
var instance *Container
var once sync.Once

// GetContainer returns the singleton container
func GetContainer() *Container {
	once.Do(func() {
		instance = NewContainer()
	})
	return instance
}

// registerDefaultServices registers default services
func registerDefaultServices(c *Container) {
	// Register HTTP client
	httpClient := utils.NewHTTPClient()
	c.Register("httpClient", httpClient)

	// Register cache
	cacheInstance := cache.NewCache("", 24*time.Hour, true, false)
	c.Register("cache", cacheInstance)

	// Register repository
	wordRepo := repository.NewWordRepository()
	c.Register("wordRepository", wordRepo)

	// Register DudenScraper
	dudenScraper := crawler.NewDudenScraper()
	c.Register("dudenScraper", dudenScraper)

	// Register cached scraper
	cachedDudenScraper := crawler.NewCachedDudenScraper(dudenScraper, cacheInstance)
	c.Register("cachedDudenScraper", cachedDudenScraper)

	// Register WordService
	wordService := services.NewWordService(cachedDudenScraper, wordRepo)
	c.Register("wordService", wordService)

	// Register BatchService
	batchService := services.NewBatchService(wordService, wordRepo, 5, 30*time.Second)
	c.Register("batchService", batchService)
}

// Helper functions for getting services

// GetWordService returns the WordService from the singleton container
func GetWordService() *services.WordService {
	return GetContainer().GetWordService()
}

// GetDudenScraper returns the DudenScraper from the singleton container
func GetDudenScraper() *crawler.DudenScraper {
	return GetContainer().GetDudenScraper()
}

// GetCachedDudenScraper returns the CachedDudenScraper from the singleton container
func GetCachedDudenScraper() *crawler.CachedDudenScraper {
	return GetContainer().GetCachedDudenScraper()
}

// GetCache returns the Cache from the singleton container
func GetCache() *cache.Cache {
	return GetContainer().GetCache()
}

// GetHTTPClient returns the HTTPClient from the singleton container
func GetHTTPClient() *utils.HTTPClient {
	return GetContainer().MustGet("httpClient").(*utils.HTTPClient)
}

// GetWordRepository returns the WordRepository from the singleton container
func GetWordRepository() *repository.WordRepository {
	return GetContainer().GetWordRepository()
}

// GetBatchService returns the BatchService from the singleton container
func GetBatchService() *services.BatchService {
	return GetContainer().GetBatchService()
}
