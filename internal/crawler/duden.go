// ./internal/crawler/duden.go
package crawler

import (
	"fmt"

	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
)

// CrawlDuden is a placeholder function for now
func CrawlDuden(word string) (*models.Word, error) {
	fmt.Println("🕵️ Crawling word:", word)

	// Placeholder return value
	return &models.Word{
		Word:     word,
		Meanings: []string{"Definition placeholder"},
		Synonyms: []string{"Synonym placeholder"},
		Grammar:  "Grammar placeholder",
	}, nil
}
