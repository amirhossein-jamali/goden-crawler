// File: internal/crawler/duden_scraper.go

package crawler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/amirhossein-jamali/goden-crawler/internal/crawler/extractors"
	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// DudenScraper scrapes data from the Duden website
// It implements the DudenCrawler interface
type DudenScraper struct {
	client           *http.Client
	extractorFactory *extractors.ExtractorFactory
	baseURL          string
	searchURL        string
	headers          map[string]string
}

// NewDudenScraper creates a new DudenScraper
func NewDudenScraper() *DudenScraper {
	return &DudenScraper{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		extractorFactory: extractors.NewExtractorFactory(),
		baseURL:          "https://www.duden.de",
		searchURL:        "https://www.duden.de/suchen/dudenonline/",
		headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
	}
}

// FetchWordData fetches data for a word
func (s *DudenScraper) FetchWordData(word string) (map[string]string, error) {
	doc, err := s.fetchWordDoc(word)
	if err != nil {
		return nil, err
	}

	// Get available sections
	sections := s.GetAvailableSections()

	// Extract data from each section
	data := make(map[string]string)
	for _, section := range sections {
		extractor, err := s.extractorFactory.CreateExtractor(section, doc)
		if err != nil {
			fmt.Printf("⚠️ Failed to create extractor for section '%s': %v\n", section, err)
			continue
		}

		extractedData := extractor.Extract()
		if extractedData != nil {
			// Convert the extracted data to a string representation
			data[section] = fmt.Sprintf("%v", extractedData)
		}
	}

	return data, nil
}

// FetchWordDataStructured fetches data for a word and returns a structured Word object
func (s *DudenScraper) FetchWordDataStructured(word string) (*models.Word, error) {
	doc, err := s.fetchWordDoc(word)
	if err != nil {
		return nil, err
	}

	// Create a Word object
	wordData := &models.Word{}

	// Extract general info
	generalInfoExtractor, _ := s.extractorFactory.CreateExtractor("general_info", doc)
	if generalInfo, ok := generalInfoExtractor.Extract().(map[string]interface{}); ok {
		wordData.Word = fmt.Sprintf("%v", generalInfo["word"])
		wordData.Article = fmt.Sprintf("%v", generalInfo["article"])

		if wordTypes, ok := generalInfo["word_type"].([]string); ok {
			wordData.WordType = wordTypes
		}

		wordData.Frequency = fmt.Sprintf("%v", generalInfo["frequency"])

		if pronunciations, ok := generalInfo["pronunciation"].([]models.Pronunciation); ok {
			wordData.Pronunciation = pronunciations
		}
	}

	// Extract meanings
	meaningExtractor, _ := s.extractorFactory.CreateExtractor("bedeutungen", doc)
	if meanings, ok := meaningExtractor.Extract().([]models.Meaning); ok {
		wordData.Meanings = meanings
	}

	// Extract synonyms
	synonymExtractor, _ := s.extractorFactory.CreateExtractor("synonyme", doc)
	if synonymData, ok := synonymExtractor.Extract().(map[string]interface{}); ok {
		if synonyms, ok := synonymData["synonyms"].([]models.Synonym); ok {
			wordData.Synonyms = synonyms
		}
	}

	// Extract grammar
	grammarExtractor, _ := s.extractorFactory.CreateExtractor("grammatik", doc)
	if grammarData, ok := grammarExtractor.Extract().(map[string]interface{}); ok {
		if paragraphs, ok := grammarData["paragraphs"].(map[string]interface{}); ok {
			if text, ok := paragraphs["text"].(string); ok {
				wordData.Grammar = text
			}
		}
	}

	// Extract origin
	originExtractor, _ := s.extractorFactory.CreateExtractor("herkunft", doc)
	if origins, ok := originExtractor.Extract().([]models.Origin); ok {
		wordData.Origin = origins
	}

	// Extract fun facts
	funFactsExtractor, _ := s.extractorFactory.CreateExtractor("wussten_sie_schon", doc)
	if funFacts, ok := funFactsExtractor.Extract().([]string); ok {
		wordData.FunFacts = funFacts
	}

	return wordData, nil
}

// GetSuggestions gets suggestions for a word
func (s *DudenScraper) GetSuggestions(word string) ([]models.Synonym, error) {
	encodedWord := url.QueryEscape(word)
	searchURL := fmt.Sprintf("%s%s", s.searchURL, encodedWord)

	doc, err := s.makeRequest(searchURL)
	if err != nil {
		return nil, err
	}

	var suggestions []models.Synonym
	doc.Find("a[href^='/rechtschreibung/'], a[href^='/node/'], a[href^='/fremdwort/']").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		href, exists := sel.Attr("href")
		if exists && text != "" {
			suggestions = append(suggestions, models.Synonym{
				Text: text,
				Link: s.baseURL + href,
			})
		}
	})

	return suggestions, nil
}

// GetAvailableSections returns a list of all available sections
func (s *DudenScraper) GetAvailableSections() []string {
	return s.extractorFactory.GetAvailableSections()
}

// fetchWordDoc fetches the HTML document for a word
func (s *DudenScraper) fetchWordDoc(word string) (*goquery.Document, error) {
	encodedWord := url.QueryEscape(word)
	wordURL := fmt.Sprintf("%s/rechtschreibung/%s", s.baseURL, encodedWord)

	// Try direct URL first
	doc, err := s.makeRequest(wordURL)
	if err == nil {
		// Check if it's an error page
		title := doc.Find("title").Text()
		if !strings.Contains(title, "Fehlermeldung") {
			return doc, nil
		}
	}

	// If direct URL fails, try to find suggestions
	fmt.Printf("Word '%s' not found. Searching for alternatives...\n", word)
	suggestions, err := s.GetSuggestions(word)
	if err != nil || len(suggestions) == 0 {
		return nil, fmt.Errorf("no alternatives found for '%s'", word)
	}

	// Try each suggestion
	for _, suggestion := range suggestions {
		fmt.Printf("Trying alternative: %s (%s)\n", suggestion.Text, suggestion.Link)

		// Add a small delay to avoid rate limiting
		time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)

		doc, err := s.makeRequest(suggestion.Link)
		if err == nil {
			return doc, nil
		}
	}

	return nil, errors.New("failed to find any valid alternatives")
}

// makeRequest makes an HTTP request and returns a goquery document
func (s *DudenScraper) makeRequest(url string) (*goquery.Document, error) {
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range s.headers {
		req.Header.Add(key, value)
	}

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("page not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
