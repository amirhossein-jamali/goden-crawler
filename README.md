# Goden Crawler

## 📖 Overview
**Goden Crawler** is a powerful CLI tool written in **Golang** that scrapes linguistic data from the **Duden** online dictionary. It allows users to extract information such as word meanings, synonyms, grammar details, and more. The project is designed to be modular, scalable, and easy to extend.

## 🚀 Features
- **Deep crawling**: Extracts all related linguistic information from Duden.
- **CLI interface**: Easily fetch data from the command line.
- **Format selection**: Supports both **JSON** and **text** output.
- **Modular design**: Extensible architecture with separate extractors for different data types.
- **Rate limiting**: Prevents being blocked by Duden using request delays.

## 📂 Project Structure
```
goden-crawler/
│── go.mod                   # Go module file
│── go.sum                   # Dependency lock file
│── README.md                # Project documentation
│── cmd/                     # CLI commands
│   ├── root.go              # Root CLI command
│   ├── scrape.go            # CLI command for scraping a word
│── internal/                # Core logic
│   ├── crawler/             # Web scraping logic
│   │   ├── duden.go         # Scraper for Duden
│   ├── extractors/          # Data extractors for different sections
│   │   ├── extractor.go     # Base extractor logic
│   │   ├── general.go       # Extractor for general info
│   │   ├── meanings.go      # Extractor for meanings
│   │   ├── grammar.go       # Extractor for grammar info
│   │   ├── synonyms.go      # Extractor for synonyms
│── pkg/                     # Reusable packages
│   ├── models/              # Data models
│   │   ├── word.go          # Model for storing word data
│   ├── formatter/           # Formatting logic
│   │   ├── json.go          # JSON output formatter
│   │   ├── text.go          # Text output formatter
│── tests/                   # Unit and integration tests
│   ├── test_crawler.go      # Tests for web scraping logic
│── configs/                 # Configuration files (if needed)
│── main.go                  # Entry point for CLI
```

## 🛠 Installation
### Prerequisites
Ensure you have **Go 1.20+** installed. Then, clone the repository and initialize the project:

```sh
git clone https://github.com/amirhossein-jamali/goden-crawler.git
cd goden-crawler
go mod tidy
```

## 🔧 Usage
### Fetch Data for a Word
To scrape data from Duden for a given word:
```sh
go run main.go scrape <word>
```
Example:
```sh
go run main.go scrape Lernen
```

### Choose Output Format
```sh
go run main.go scrape Lernen -o json
```
Supported formats: **text** (default) and **json**.

## 🧪 Running Tests
To run tests:
```sh
go test ./...
```

## 🚀 Contributing
1. Fork the repository.
2. Create a new branch (`feature-new`).
3. Commit your changes.
4. Push the branch and submit a pull request.

## 📄 License
This project is licensed under the **MIT License**.