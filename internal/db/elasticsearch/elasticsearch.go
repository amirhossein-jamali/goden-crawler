package elasticsearch

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var esClient *elasticsearch.Client

// ConnectElasticsearch establishes a connection to Elasticsearch
func ConnectElasticsearch() *elasticsearch.Client {
	if esClient != nil {
		return esClient
	}

	// Get Elasticsearch URI from environment variable or use default
	esURI := os.Getenv("ELASTICSEARCH_URI")
	if esURI == "" {
		esURI = "http://localhost:9200"
	}

	cfg := elasticsearch.Config{
		Addresses: []string{esURI},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	// Test the connection
	res, err := client.Info()
	if err != nil {
		log.Fatalf("Failed to get Elasticsearch info: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Elasticsearch error: %s", res.String())
	}

	esClient = client
	return client
}

// CreateWordsIndex creates the index for words if it doesn't exist
func CreateWordsIndex() error {
	client := ConnectElasticsearch()

	// Check if index exists
	res, err := client.Indices.Exists([]string{"words"})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// If index exists, return
	if res.StatusCode == 200 {
		return nil
	}

	// Define index mapping
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"word": {
					"type": "text",
					"analyzer": "standard",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				},
				"partOfSpeech": {
					"type": "keyword"
				},
				"definitions": {
					"type": "text",
					"analyzer": "standard"
				},
				"examples": {
					"type": "text",
					"analyzer": "standard"
				},
				"synonyms": {
					"type": "nested",
					"properties": {
						"text": {
							"type": "text",
							"analyzer": "standard"
						}
					}
				}
			}
		}
	}`

	// Create index
	res, err = client.Indices.Create(
		"words",
		client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}

	return nil
}

// IndexWord indexes a word in Elasticsearch
func IndexWord(word *models.Word) error {
	client := ConnectElasticsearch()

	// Ensure index exists
	if err := CreateWordsIndex(); err != nil {
		return err
	}

	// Convert word to JSON
	wordJSON, err := json.Marshal(word)
	if err != nil {
		return err
	}

	// Index document
	req := esapi.IndexRequest{
		Index:      "words",
		DocumentID: word.Word,
		Body:       strings.NewReader(string(wordJSON)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}

	return nil
}

// SearchWords searches for words in Elasticsearch
func SearchWords(query string) ([]models.Word, error) {
	client := ConnectElasticsearch()

	// Build search query
	searchQuery := `{
		"query": {
			"multi_match": {
				"query": "` + query + `",
				"fields": ["word^3", "definitions", "examples"]
			}
		}
	}`

	// Perform search
	res, err := client.Search(
		client.Search.WithIndex("words"),
		client.Search.WithBody(strings.NewReader(searchQuery)),
		client.Search.WithSize(10),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, err
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract hits
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	words := make([]models.Word, 0, len(hits))

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceJSON, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var word models.Word
		if err := json.Unmarshal(sourceJSON, &word); err != nil {
			continue
		}

		words = append(words, word)
	}

	return words, nil
}
