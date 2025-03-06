package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	ctx    context.Context
)

// ConnectMongoDB establishes a connection to MongoDB
func ConnectMongoDB() *mongo.Client {
	if client != nil {
		return client
	}

	// Get MongoDB URI from environment variable or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/goden-crawler"
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
	return client
}

// GetWordsCollection returns the words collection
func GetWordsCollection() *mongo.Collection {
	client := ConnectMongoDB()
	return client.Database("goden-crawler").Collection("words")
}

// SaveWord saves a word to MongoDB
func SaveWord(word *models.Word) error {
	collection := GetWordsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use upsert to update if exists or insert if not
	filter := bson.M{"word": word.Word}
	update := bson.M{"$set": word}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetWord retrieves a word from MongoDB
func GetWord(wordText string) (*models.Word, error) {
	collection := GetWordsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var word models.Word
	filter := bson.M{"word": wordText}
	err := collection.FindOne(ctx, filter).Decode(&word)
	if err != nil {
		return nil, err
	}

	return &word, nil
}

// Close closes the MongoDB connection
func Close() error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return client.Disconnect(ctx)
	}
	return nil
}
