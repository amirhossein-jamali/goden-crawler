package redis

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// ConnectRedis establishes a connection to Redis
func ConnectRedis() *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	// Get Redis URI from environment variable or use default
	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		redisURI = "redis://localhost:6379/0"
	}

	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		log.Fatalf("Failed to parse Redis URI: %v", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	redisClient = client
	return client
}

// CacheWord caches a word in Redis
func CacheWord(word *models.Word, ttl time.Duration) error {
	client := ConnectRedis()
	ctx := context.Background()

	// Convert word to JSON
	data, err := json.Marshal(word)
	if err != nil {
		return err
	}

	// Cache with TTL
	key := "word:" + word.Word
	return client.Set(ctx, key, data, ttl).Err()
}

// GetCachedWord retrieves a cached word from Redis
func GetCachedWord(wordText string) (*models.Word, error) {
	client := ConnectRedis()
	ctx := context.Background()

	// Get from cache
	key := "word:" + wordText
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var word models.Word
	if err := json.Unmarshal(data, &word); err != nil {
		return nil, err
	}

	return &word, nil
}

// DeleteCachedWord removes a word from the cache
func DeleteCachedWord(wordText string) error {
	client := ConnectRedis()
	ctx := context.Background()

	key := "word:" + wordText
	return client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func Close() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
