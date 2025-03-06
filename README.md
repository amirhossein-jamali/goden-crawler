# 📚 Goden Crawler

Goden Crawler is a powerful command-line interface (CLI) tool designed to extract detailed linguistic information from the **Duden** online dictionary. Built with **Golang** and **Cobra** for CLI management, it features a modular and scalable architecture following clean architecture principles, making it easy to maintain and extend.

---

## 🚀 Features

- **Comprehensive Information Extraction**:
  - Word details (article, type, frequency)
  - Grammatical information (conjugations, declensions)
  - Meanings and definitions with examples
  - Synonyms and related words
  - Etymology and word origins
  - Pronunciation guides
  - Spelling information
  - "Did you know?" fun facts

- **Multiple Operation Modes**:
  - Single word scraping
  - Interactive shell mode
  - Batch processing with concurrent workers
  - Suggestions for similar words

- **Flexible Output Formats**:
  - Text (human-readable)
  - JSON (machine-readable)
  - Extensible formatter system

- **Performance Optimizations**:
  - Caching system (memory and disk)
  - Concurrent processing
  - Configurable timeouts and retries

- **Developer-Friendly**:
  - Clean architecture (domain, application, infrastructure layers)
  - Dependency injection
  - Extensive logging
  - Error handling
  - Shell completion support

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

2. **Install dependencies**:

```bash
go mod tidy
```

3. **Build the application**:

```bash
go build -o goden-crawler
```

4. **Run the application**:

```bash
./goden-crawler --help
```

---

## 🚀 Usage

### Single Word Scraping

Extract linguistic information for a single German word:

```bash
./goden-crawler scrape <word> [--output format]
```

Options:
- `--output`, `-o`: Output format (text, json). Default: text

Example:
```bash
./goden-crawler scrape zahlen --output json
```

### Interactive Mode

Start an interactive shell for continuous word lookups:

```bash
./goden-crawler interactive
```

In interactive mode, you can:
- Type a word to get its information
- Use `suggest <word>` to get word suggestions
- Switch output format with `format json` or `format text`
- View available sections with `sections`
- Get help with `help`
- Exit with `exit` or `quit`

### Batch Processing

Process multiple words concurrently and save results to files:

```bash
./goden-crawler batch <word1> <word2> <word3> ... [flags]
```

Options:
- `--output`, `-o`: Output format (text, json). Default: json
- `--workers`, `-w`: Number of concurrent workers. Default: 5
- `--timeout`, `-t`: Timeout in seconds per word. Default: 30
- `--prefix`, `-p`: Output filename prefix. Default: none

Example:
```bash
./goden-crawler batch haus baum auto --workers 3 --output json --prefix "german_"
```

This will create files: german_haus.json, german_baum.json, german_auto.json

### Shell Completion

Generate shell completion scripts:

```bash
./goden-crawler completion [bash|zsh|fish|powershell]
```

---

## 📊 Sample Output

### Text Format

```
--- GENERAL INFO ---
Word: zahlen
Article: 
Word Type: Verb
Frequency: ★★★★☆ (high)

--- GRAMMAR ---
Conjugation: zahlt, zahlte, hat gezahlt

--- MEANINGS ---
1. einen Geldbetrag als Gegenleistung für etwas geben
   Examples:
   - bar, mit Kreditkarte zahlen
   - die Rechnung, Miete, Steuern zahlen
   - er hat für alle gezahlt

2. einen bestimmten Preis haben
   Examples:
   - für das Haus musste er viel zahlen
   - was zahlt man für ein Kilo Äpfel?

--- SYNONYMS ---
begleichen, berappen, bestreiten, bezahlen, entrichten, erstatten...
```

### JSON Format

```json
{
  "word": "zahlen",
  "word_type": ["Verb"],
  "frequency": "★★★★☆",
  "grammar": "zahlt, zahlte, hat gezahlt",
  "meanings": [
    {
      "text": "einen Geldbetrag als Gegenleistung für etwas geben",
      "examples": [
        "bar, mit Kreditkarte zahlen",
        "die Rechnung, Miete, Steuern zahlen",
        "er hat für alle gezahlt"
      ]
    },
    {
      "text": "einen bestimmten Preis haben",
      "examples": [
        "für das Haus musste er viel zahlen",
        "was zahlt man für ein Kilo Äpfel?"
      ]
    }
  ],
  "synonyms": [
    {"text": "begleichen"},
    {"text": "berappen"},
    {"text": "bestreiten"},
    {"text": "bezahlen"},
    {"text": "entrichten"},
    {"text": "erstatten"}
  ]
}
```

---

## 🧪 Testing

Run tests to ensure everything works as expected:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/crawler/...
```

---

## 📂 Project Structure

The project follows clean architecture principles with clear separation of concerns:

```
goden-crawler/
├── cmd/                     # CLI commands using Cobra
│   ├── root.go              # Main entry point for CLI
│   ├── scrape.go            # Single word scraping
│   ├── interactive.go       # Interactive shell mode
│   ├── batch.go             # Batch processing
│   └── completion.go        # Shell completion
├── internal/                # Internal packages (not importable)
│   ├── domain/              # Core domain models and interfaces
│   │   ├── models/          # Domain models
│   │   │   └── word_builder.go # Builder for word models
│   │   └── interfaces/      # Core interfaces
│   │       ├── crawler.go   # Crawler interface
│   │       ├── extractor.go # Extractor interface
│   │       └── service.go   # Service interfaces
│   ├── application/         # Application services
│   │   └── services/        # Business logic services
│   │       ├── word_service.go  # Word data operations
│   │       └── batch_service.go # Batch processing
│   ├── infrastructure/      # External services implementation
│   │   ├── crawler/         # Web scraping implementation
│   │   │   ├── duden_scraper.go     # Duden website scraper
│   │   │   └── cached_duden_scraper.go # Cached scraper
│   │   ├── extractors/      # Data extraction modules
│   │   │   ├── base.go      # Base extractor
│   │   │   ├── factory.go   # Extractor factory
│   │   │   ├── strategy.go  # Extraction strategy pattern
│   │   │   └── general_info.go # General info extractor
│   │   ├── cache/           # Caching implementation
│   │   │   └── cache.go     # Memory and disk caching
│   │   ├── http/            # HTTP client implementation
│   │   │   └── client.go    # Custom HTTP client
│   │   ├── middleware/      # Middleware chain
│   │   │   └── chain.go     # Middleware implementation
│   │   ├── container/       # Dependency injection
│   │   │   └── container.go # Service container
│   │   ├── events/          # Event system
│   │   │   └── observer.go  # Observer pattern
│   │   └── plugins/         # Plugin system
│   │       ├── plugin.go    # Plugin manager
│   │       └── plugin_example.go # Example plugin
│   ├── crawler/             # Crawler implementation
│   │   ├── duden.go         # Duden crawler interface
│   │   └── extractors/      # Specific extractors
│   │       ├── base_extractor.go      # Base extractor
│   │       ├── bedeutungen_extractor.go # Meanings extractor
│   │       ├── grammatik_extractor.go   # Grammar extractor
│   │       ├── synonyme_extractor.go    # Synonyms extractor
│   │       ├── herkunft_extractor.go    # Origin extractor
│   │       ├── rechtschreibung_extractor.go # Spelling extractor
│   │       └── wussten_sie_schon_extractor.go # Fun facts extractor
│   └── formatter/           # Output formatting
│       └── formatter.go     # Text/JSON formatter
├── pkg/                     # Public packages (importable)
│   ├── models/              # Data models
│   │   └── word.go          # Word model
│   ├── utils/               # Utility functions
│   │   ├── http_client.go   # HTTP client utilities
│   │   ├── string_utils.go  # String manipulation utilities
│   │   └── config.go        # Configuration utilities
│   ├── logger/              # Logging utilities
│   │   └── logger.go        # Logger implementation
│   └── errors/              # Error handling
│       └── errors.go        # Custom errors
├── main.go                  # Entry point
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies lock file
└── README.md                # Documentation
```

---

## 🔧 Architecture

Goden Crawler follows clean architecture principles with three main layers:

1. **Domain Layer** (internal/domain)
   - Contains core business logic and interfaces
   - Independent of external frameworks and tools
   - Defines the contracts that other layers must implement

2. **Application Layer** (internal/application)
   - Implements use cases using domain interfaces
   - Orchestrates the flow of data between domain and infrastructure
   - Contains business rules specific to the application

3. **Infrastructure Layer** (internal/infrastructure)
   - Implements interfaces defined in the domain layer
   - Handles external concerns like HTTP, caching, and data persistence
   - Adapts external libraries and frameworks to the application

This architecture ensures:
- Separation of concerns
- Testability
- Maintainability
- Flexibility to change external dependencies

---

## 🧩 Design Patterns

The project implements several design patterns:

- **Builder Pattern**: For constructing Word objects (word_builder.go)
- **Factory Pattern**: For creating extractors (extractor_factory.go)
- **Strategy Pattern**: For different extraction strategies (strategy.go)
- **Observer Pattern**: For event handling (observer.go)
- **Middleware Pattern**: For request processing (chain.go)
- **Dependency Injection**: For service management (container.go)
- **Singleton Pattern**: For global access to services

---

## 📊 Technologies & Tools

This project leverages the following tools and libraries:

- **Golang 1.21+**: Core programming language
- **Cobra**: CLI framework for handling commands
- **goquery**: Web scraping library for Go
- **Go testing framework**: For unit and integration testing

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

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

### Contribution Guidelines

- Follow Go coding standards
- Write tests for new features
- Update documentation as needed
- Keep pull requests focused on a single feature/fix

---

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](https://github.com/amirhossein-jamali/goden-crawler/blob/main/LICENSE) file for details.

---

## 🙏 Acknowledgements

- [Duden](https://www.duden.de/) for providing the linguistic data
- The Go community for excellent libraries and tools
- All contributors who have helped improve this project

---

## 📂 .gitignore

To ensure that unnecessary files are not tracked by Git, create a `.gitignore` file in the root of your project with the following content:

```
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, build with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below if you use dep)
# vendor/

# Go workspace file
*.code-workspace

# IDE specific files
.idea/
.vscode/

# Logs
*.log

# Temporary files
*.tmp

# Build output
/goden-crawler

# Go module files
/go.sum

# Environment variables
.env
```

This will help keep your repository clean by ignoring files that are not necessary to track.