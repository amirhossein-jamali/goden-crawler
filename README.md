# 📚 Goden Crawler

Goden Crawler is a command-line interface (CLI) tool designed to extract detailed linguistic information from the **Duden** online dictionary. Built with **Golang** and **Cobra** for CLI management, it features a modular and scalable architecture, making it easy to maintain and extend.

---

## 🚀 Features

- **Comprehensive Information Extraction**:
    - Automatically detect word type, retrieve articles, and word frequency.
    - Fetch verb conjugations, noun declensions, and more.
    - Retrieve word meanings, usage examples, and contextual information.
- **Automated Testing**: Includes unit and integration tests using **Go testing framework**.
- **User-Friendly CLI**: Easy-to-use command-line interface for quick word analysis.

---

## 📦 Prerequisites

Before getting started, make sure you have the following installed:

- **Golang 1.21+**
- **Git** (for version control)

---

## ⚙️ Project Setup

Follow these steps to set up the project:

1. **Clone the repository**:

```bash
git clone https://github.com/amirhossein-jamali/goden-crawler.git
cd goden-crawler
```

2. **Initialize Go modules** (if not already initialized):

```bash
go mod tidy
```

---

## 🚀 Usage

Run the tool to extract linguistic information from Duden.

### General Usage:

```bash
go run main.go scrape <word>
```

- `<word>`: The German word to analyze (word type is automatically detected).

### Examples:

#### Analyze a Word:

```bash
go run main.go scrape zahlen
```

**Sample Output**:

```
--- GENERAL INFO ---
{
    "word": "zahlen",
    "article": "das",
    "word_type": "Verb"
}

--- GRAMMAR ---
{
    "conjugation": "zahlt, zahlte, hat gezahlt"
}
```

---

## 🧪 Testing

Run tests to ensure everything works as expected:

- Run all tests:

```bash
go test ./...
```

- Run specific test files:

```bash
go test ./tests/scraper_test.go
```

---

## 📂 Project Structure

```plaintext
goden-crawler/
├── cmd/                     # CLI commands using Cobra
│   ├── root.go              # Main entry point for CLI
│   ├── scrape.go            # Handles scraping words from Duden
├── internal/
│   ├── crawler/             # Core web scraper
│   │   ├── duden_scraper.go # Fetches data from Duden website
│   ├── extractors/          # Extraction logic for different sections
│   │   ├── base.go          # Base extractor interface
│   │   ├── general.go       # Extracts general word info
│   │   ├── meanings.go      # Extracts meanings and examples
│   │   ├── grammar.go       # Extracts grammatical information
│   │   ├── synonyms.go      # Extracts synonyms
│   │   ├── spelling.go      # Extracts spelling information
│   │   ├── origin.go        # Extracts word origins
│   │   ├── funfact.go       # Extracts "Did you know?" section
│   │   ├── extractor_factory.go  # Factory for handling extractors
│   ├── formatter/           # Formatting output
│   │   ├── formatter.go     # Converts extracted data to JSON/Text
│   ├── models/              # Data models
│   │   ├── word.go          # Defines the `Word` struct
├── tests/                   # Unit & Integration tests
│   ├── scraper_test.go      # Tests for DudenScraper
│   ├── extractors_test.go   # Tests for extractors
├── main.go                  # Entry point
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies lock file
├── README.md                # Documentation
```

---

## 📊 Technologies & Tools

This project leverages the following tools and libraries:

- **Golang 1.21+**: Core programming language.
- **Cobra**: CLI framework for handling commands.
- **Colly**: Web scraping library for Go.
- **Go testing framework**: For unit and integration testing.

---

## 🤝 Contributing

We welcome contributions! Here’s how you can help:

1. **Fork the repository**: 🍴
2. **Create a new branch**: 🌿

```bash
git checkout -b feature/your-feature-name
```

3. **Make your changes**: 💡
4. **Commit your changes**: 🖍️

```bash
git commit -m "Add your feature"
```

5. **Push to your branch**: 🚀

```bash
git push origin feature/your-feature-name
```

6. **Open a Pull Request**: 🔥

---

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](https://github.com/amirhossein-jamali/goden-crawler/blob/main/LICENSE) file for details.

---