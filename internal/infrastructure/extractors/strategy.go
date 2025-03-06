// File: internal/infrastructure/extractors/strategy.go

package extractors

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
)

// ExtractionStrategy defines the interface for extraction strategies
type ExtractionStrategy interface {
	// Extract extracts data from the document
	Extract(doc *goquery.Document) (interface{}, error)

	// GetName returns the name of the strategy
	GetName() string
}

// ExtractionContext holds the current extraction strategy and document
type ExtractionContext struct {
	strategy ExtractionStrategy
	doc      *goquery.Document
}

// NewExtractionContext creates a new extraction context
func NewExtractionContext(doc *goquery.Document) *ExtractionContext {
	return &ExtractionContext{
		doc: doc,
	}
}

// SetStrategy sets the extraction strategy
func (c *ExtractionContext) SetStrategy(strategy ExtractionStrategy) {
	c.strategy = strategy
	logger.Debug("Set extraction strategy", logger.F("strategy", strategy.GetName()))
}

// ExecuteStrategy executes the current strategy
func (c *ExtractionContext) ExecuteStrategy() (interface{}, error) {
	if c.strategy == nil {
		return nil, errors.New("no extraction strategy set")
	}

	if c.doc == nil {
		return nil, errors.New("no document set")
	}

	logger.Debug("Executing extraction strategy", logger.F("strategy", c.strategy.GetName()))
	return c.strategy.Extract(c.doc)
}

// Global strategy registry
var globalRegistry *StrategyRegistry

// StrategyRegistry holds all registered extraction strategies
type StrategyRegistry struct {
	strategies map[string]ExtractionStrategy
}

// NewStrategyRegistry creates a new strategy registry
func NewStrategyRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		strategies: make(map[string]ExtractionStrategy),
	}
}

// Register registers a strategy
func (r *StrategyRegistry) Register(strategy ExtractionStrategy) {
	r.strategies[strategy.GetName()] = strategy
	logger.Debug("Registered extraction strategy", logger.F("strategy", strategy.GetName()))
}

// Get retrieves a strategy by name
func (r *StrategyRegistry) Get(name string) (ExtractionStrategy, bool) {
	strategy, exists := r.strategies[name]
	return strategy, exists
}

// GetAll returns all registered strategies
func (r *StrategyRegistry) GetAll() map[string]ExtractionStrategy {
	// Create a copy to avoid concurrent map access
	strategies := make(map[string]ExtractionStrategy, len(r.strategies))
	for name, strategy := range r.strategies {
		strategies[name] = strategy
	}
	return strategies
}

// Global registry functions

// RegisterStrategy registers a strategy with the global registry
func RegisterStrategy(strategy ExtractionStrategy) {
	if globalRegistry == nil {
		globalRegistry = NewStrategyRegistry()
	}
	globalRegistry.Register(strategy)
}

// GetStrategy retrieves a strategy from the global registry
func GetStrategy(name string) (ExtractionStrategy, bool) {
	if globalRegistry == nil {
		globalRegistry = NewStrategyRegistry()
	}
	return globalRegistry.Get(name)
}

// GetAllStrategies returns all registered strategies from the global registry
func GetAllStrategies() map[string]ExtractionStrategy {
	if globalRegistry == nil {
		globalRegistry = NewStrategyRegistry()
	}
	return globalRegistry.GetAll()
}
