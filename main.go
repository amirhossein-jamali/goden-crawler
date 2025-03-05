// ./main.go
package main

import (
	"fmt"
	
	"github.com/amirhossein-jamali/goden-crawler/internal/crawler"
)

func main() {
	scraper := crawler.NewDudenScraper()
	word := "Lernen"

	fmt.Printf("🔍 Fetching data for: %s\n", word)
	data, err := scraper.FetchWordData(word)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println(data)
}
