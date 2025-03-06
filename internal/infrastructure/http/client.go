// File: internal/infrastructure/http/client.go

package http

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	customErrors "github.com/amirhossein-jamali/goden-crawler/pkg/errors"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
)

// Client defines the interface for HTTP clients
type Client interface {
	// Get makes a GET request to the specified URL
	Get(url string) (*http.Response, error)

	// GetWithHeaders makes a GET request with custom headers
	GetWithHeaders(url string, headers map[string]string) (*http.Response, error)

	// GetBody makes a GET request and returns the response body as a string
	GetBody(url string) (string, error)
}

// DefaultClient is the default implementation of the Client interface
type DefaultClient struct {
	client  *http.Client
	headers map[string]string
	retries int
}

// NewClient creates a new DefaultClient with default settings
func NewClient() *DefaultClient {
	return &DefaultClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
		retries: 3,
	}
}

// WithTimeout sets the timeout for the client
func (c *DefaultClient) WithTimeout(timeout time.Duration) *DefaultClient {
	c.client.Timeout = timeout
	return c
}

// WithHeaders sets the headers for the client
func (c *DefaultClient) WithHeaders(headers map[string]string) *DefaultClient {
	c.headers = headers
	return c
}

// WithRetries sets the number of retries for the client
func (c *DefaultClient) WithRetries(retries int) *DefaultClient {
	c.retries = retries
	return c
}

// Get makes a GET request to the specified URL
func (c *DefaultClient) Get(url string) (*http.Response, error) {
	return c.GetWithHeaders(url, nil)
}

// GetWithHeaders makes a GET request with custom headers
func (c *DefaultClient) GetWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add default headers
	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	// Add custom headers
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	// Add a random delay to avoid rate limiting
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	// Try to make the request with retries
	var resp *http.Response
	var lastErr error

	for i := 0; i < c.retries; i++ {
		resp, err = c.client.Do(req)
		if err == nil {
			// Check if the response is successful
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return resp, nil
			}

			// Close the response body to avoid resource leaks
			resp.Body.Close()

			// Create an error for the unsuccessful status code
			lastErr = fmt.Errorf("unsuccessful status code: %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		logger.Warn("Request failed",
			logger.F("attempt", i+1),
			logger.F("retries", c.retries),
			logger.F("error", err))

		// If this is not the last retry, wait before trying again
		if i < c.retries-1 {
			// Exponential backoff
			backoff := time.Duration(1<<i) * time.Second
			logger.Info("Retrying", logger.F("backoff", backoff))
			time.Sleep(backoff)
		}
	}

	// If we get here, all retries failed
	if lastErr != nil {
		return nil, customErrors.Wrap(lastErr, "all retries failed")
	}

	// This should never happen, but just in case
	return nil, fmt.Errorf("all retries failed with unknown error")
}

// GetDocument makes a GET request and returns the response as a goquery document
func (c *DefaultClient) GetDocument(url string) (*goquery.Document, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// GetBody makes a GET request and returns the response body as a string
func (c *DefaultClient) GetBody(url string) (string, error) {
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
