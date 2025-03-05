// ./internal/crawler/duden_scraper.go

package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type DudenScraper struct {}

// NewDudenScraper initializes a new scraper instance
func NewDudenScraper() *DudenScraper {
	return &DudenScraper{}
}

// FetchWordData scrapes the Duden website for a given word
func (ds *DudenScraper) FetchWordData(word string) (string, error) {
	word = strings.ToLower(word) // Convert to lowercase to match Duden URL format
	url := fmt.Sprintf("https://www.duden.de/rechtschreibung/%s", word)

	// Use net/http to fetch the page manually
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract word definition - check multiple possible selectors
	definition := doc.Find(".enumeration__text").First().Text()
	if definition == "" {
		definition = doc.Find(".entry__content p").First().Text()
	}

	if definition == "" {
		return "❌ No definition found. Please check the word or website structure.", nil
	}

	return fmt.Sprintf("✅ Definition: %s", strings.TrimSpace(definition)), nil
}