// File: pkg/utils/http_client.go

package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	customErrors "github.com/amirhossein-jamali/goden-crawler/pkg/errors"
)

// HTTPClient provides utilities for making HTTP requests
type HTTPClient struct {
	client  *http.Client
	headers map[string]string
	retries int
}

// NewHTTPClient creates a new HTTPClient with default settings
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
		retries: 3,
	}
}

// WithTimeout sets the timeout for the HTTP client
func (c *HTTPClient) WithTimeout(timeout time.Duration) *HTTPClient {
	c.client.Timeout = timeout
	return c
}

// WithHeaders sets the headers for the HTTP client
func (c *HTTPClient) WithHeaders(headers map[string]string) *HTTPClient {
	c.headers = headers
	return c
}

// WithRetries sets the number of retries for the HTTP client
func (c *HTTPClient) WithRetries(retries int) *HTTPClient {
	c.retries = retries
	return c
}

// GetDocument makes an HTTP GET request and returns a goquery document
func (c *HTTPClient) GetDocument(url string) (*goquery.Document, error) {
	var (
		resp *http.Response
		err  error
	)

	// Try the request with retries
	for i := 0; i < c.retries; i++ {
		// Create a new request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, customErrors.NewScraperError("", "failed to create request", err)
		}

		// Add headers
		for key, value := range c.headers {
			req.Header.Add(key, value)
		}

		// Make the request
		resp, err = c.client.Do(req)
		if err == nil {
			break
		}

		// If this is not the last retry, wait before trying again
		if i < c.retries-1 {
			// Exponential backoff with jitter
			backoff := time.Duration(1000*(1<<i)+rand.Intn(500)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	// If all retries failed, return the last error
	if err != nil {
		return nil, customErrors.NewScraperError("", "failed to make request after retries", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, customErrors.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, customErrors.NewScraperError("", fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, customErrors.NewScraperError("", "failed to parse HTML", err)
	}

	return doc, nil
}
