// File: internal/infrastructure/middleware/chain.go

package middleware

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
)

// Request represents a request to be processed by the middleware chain
type Request struct {
	// Context for the request
	Context context.Context

	// Data is the data to be processed
	Data interface{}
}

// Handler defines the interface for middleware handlers
type Handler interface {
	// Handle processes a request and returns a response
	Handle(request *Request) (interface{}, error)

	// SetNext sets the next handler in the chain
	SetNext(handler Handler)
}

// BaseHandler provides a base implementation of the Handler interface
type BaseHandler struct {
	next Handler
}

// SetNext sets the next handler in the chain
func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

// ProcessNext processes the next handler in the chain
func (h *BaseHandler) ProcessNext(request *Request) (interface{}, error) {
	if h.next == nil {
		return nil, nil
	}
	return h.next.Handle(request)
}

// Chain represents a middleware chain
type Chain struct {
	first Handler
	last  Handler
}

// NewChain creates a new middleware chain
func NewChain() *Chain {
	return &Chain{}
}

// Add adds a handler to the chain
func (c *Chain) Add(handler Handler) *Chain {
	if c.first == nil {
		c.first = handler
		c.last = handler
		return c
	}

	c.last.SetNext(handler)
	c.last = handler
	return c
}

// Process processes a request through the chain
func (c *Chain) Process(request *Request) (interface{}, error) {
	if c.first == nil {
		return nil, nil
	}
	return c.first.Handle(request)
}

// LoggingHandler logs requests and responses
type LoggingHandler struct {
	BaseHandler
	name string
}

// NewLoggingHandler creates a new logging handler
func NewLoggingHandler(name string) *LoggingHandler {
	return &LoggingHandler{
		name: name,
	}
}

// Handle processes a request and logs it
func (h *LoggingHandler) Handle(request *Request) (interface{}, error) {
	logger.Debug("Processing request",
		logger.F("handler", h.name),
		logger.F("data", fmt.Sprintf("%+v", request.Data)))

	response, err := h.ProcessNext(request)

	if err != nil {
		logger.Error("Error processing request",
			logger.F("handler", h.name),
			logger.F("error", err.Error()))
	} else {
		logger.Debug("Request processed successfully",
			logger.F("handler", h.name))
	}

	return response, err
}

// ValidationHandler validates requests
type ValidationHandler struct {
	BaseHandler
	validator func(interface{}) error
}

// NewValidationHandler creates a new validation handler
func NewValidationHandler(validator func(interface{}) error) *ValidationHandler {
	return &ValidationHandler{
		validator: validator,
	}
}

// Handle validates a request and processes it
func (h *ValidationHandler) Handle(request *Request) (interface{}, error) {
	if h.validator != nil {
		if err := h.validator(request.Data); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}
	return h.ProcessNext(request)
}

// CacheHandler caches responses
type CacheHandler struct {
	BaseHandler
	cache map[string]interface{}
}

// NewCacheHandler creates a new cache handler
func NewCacheHandler() *CacheHandler {
	return &CacheHandler{
		cache: make(map[string]interface{}),
	}
}

// Handle checks the cache for a response or processes the request
func (h *CacheHandler) Handle(request *Request) (interface{}, error) {
	// Generate a cache key from the request data
	cacheKey := h.generateCacheKey(request.Data)

	// Check if the response is in the cache
	if response, ok := h.cache[cacheKey]; ok {
		logger.Debug("Cache hit", logger.F("key", cacheKey))
		return response, nil
	}

	// Process the request
	response, err := h.ProcessNext(request)
	if err != nil {
		return nil, err
	}

	// Cache the response
	h.cache[cacheKey] = response
	logger.Debug("Cached response", logger.F("key", cacheKey))

	return response, nil
}

// generateCacheKey generates a cache key from the request data
func (h *CacheHandler) generateCacheKey(data interface{}) string {
	// This is a simple implementation that uses the string representation of the data
	// In a real application, you would want to use a more sophisticated approach
	dataStr := fmt.Sprintf("%+v", data)
	hash := md5.Sum([]byte(dataStr))
	return hex.EncodeToString(hash[:])
}
