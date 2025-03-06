// File: internal/infrastructure/container/container.go

package container

import (
	"fmt"
	"sync"

	"github.com/amirhossein-jamali/goden-crawler/internal/application/services"
	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
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

// Register registers a service in the container
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

// MustGet retrieves a service from the container or panics if it doesn't exist
func (c *Container) MustGet(name string) interface{} {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// GetWordService returns the word service
func (c *Container) GetWordService() *services.WordService {
	return c.MustGet("wordService").(*services.WordService)
}

// GetDudenScraper returns the Duden scraper
func (c *Container) GetDudenScraper() *crawler.DudenScraper {
	return c.MustGet("dudenScraper").(*crawler.DudenScraper)
}

// Global container instance
var globalContainer *Container
var once sync.Once

// GetContainer returns the global container instance
func GetContainer() *Container {
	once.Do(func() {
		globalContainer = NewContainer()
	})
	return globalContainer
}

// registerDefaultServices registers the default services in the container
func registerDefaultServices(c *Container) {
	// Register HTTP client
	httpClient := utils.NewHTTPClient()
	c.Register("httpClient", httpClient)

	// Register Duden scraper
	dudenScraper := crawler.NewDudenScraper()
	c.Register("dudenScraper", dudenScraper)

	// Register word service
	wordService := services.NewWordService(dudenScraper)
	c.Register("wordService", wordService)
}

// Helper functions for global container

// GetWordService returns the word service from the global container
func GetWordService() *services.WordService {
	return GetContainer().GetWordService()
}

// GetDudenScraper returns the Duden scraper from the global container
func GetDudenScraper() *crawler.DudenScraper {
	return GetContainer().GetDudenScraper()
}

// GetHTTPClient returns the HTTP client from the global container
func GetHTTPClient() *utils.HTTPClient {
	return GetContainer().MustGet("httpClient").(*utils.HTTPClient)
}
