package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/amirhossein-jamali/goden-crawler/pkg/models"
	_ "github.com/lib/pq"
)

var pgDB *sql.DB

// ConnectPostgres establishes a connection to PostgreSQL
func ConnectPostgres() *sql.DB {
	if pgDB != nil {
		return pgDB
	}

	// Get PostgreSQL URI from environment variable or use default
	pgURI := os.Getenv("POSTGRES_URI")
	if pgURI == "" {
		pgURI = "postgres://postgres:postgres@localhost:5432/goden-crawler?sslmode=disable"
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", pgURI)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	pgDB = db
	return db
}

// createTables creates the necessary tables if they don't exist
func createTables(db *sql.DB) error {
	// Create words table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id SERIAL PRIMARY KEY,
			word TEXT NOT NULL UNIQUE,
			word_type TEXT[],
			data JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create meanings table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS meanings (
			id SERIAL PRIMARY KEY,
			word_id INTEGER REFERENCES words(id) ON DELETE CASCADE,
			text TEXT NOT NULL,
			examples TEXT[],
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// SaveWord saves a word to PostgreSQL
func SaveWord(word *models.Word) error {
	db := ConnectPostgres()

	// Convert word to JSON for storage
	wordJSON, err := json.Marshal(word)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert or update word
	var wordID int
	err = tx.QueryRow(`
		INSERT INTO words (word, word_type, data)
		VALUES ($1, $2, $3)
		ON CONFLICT (word) DO UPDATE
		SET word_type = $2, data = $3, updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`, word.Word, word.WordType, wordJSON).Scan(&wordID)
	if err != nil {
		return err
	}

	// Delete existing meanings
	_, err = tx.Exec("DELETE FROM meanings WHERE word_id = $1", wordID)
	if err != nil {
		return err
	}

	// Insert meanings
	for _, meaning := range word.Meanings {
		_, err = tx.Exec(`
			INSERT INTO meanings (word_id, text, examples)
			VALUES ($1, $2, $3)
		`, wordID, meaning.Text, meaning.Examples)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

// GetWord retrieves a word from PostgreSQL
func GetWord(wordText string) (*models.Word, error) {
	db := ConnectPostgres()

	// Query for the word
	var wordJSON []byte
	err := db.QueryRow(`
		SELECT data FROM words
		WHERE word = $1
	`, wordText).Scan(&wordJSON)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var word models.Word
	if err := json.Unmarshal(wordJSON, &word); err != nil {
		return nil, err
	}

	return &word, nil
}

// Close closes the PostgreSQL connection
func Close() error {
	if pgDB != nil {
		return pgDB.Close()
	}
	return nil
}
