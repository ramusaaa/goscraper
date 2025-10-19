# GoScraper 🚀

**Enterprise-Grade Web Scraping Library & Microservice for Go**

Modern, fast, and stealth web scraping library with AI-powered extraction, anti-bot detection, and microservice architecture. Perfect for e-commerce, news, and data extraction at scale.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ramusaaa/goscraper)](https://goreportcard.com/report/github.com/ramusaaa/goscraper)
[![GoDoc](https://godoc.org/github.com/ramusaaa/goscraper?status.svg)](https://godoc.org/github.com/ramusaaa/goscraper)

## 🌟 Key Features

### 🤖 **AI-Powered Smart Extraction**
- **Multiple AI Providers**: OpenAI GPT-4, Anthropic Claude, Local models
- **Smart Content Detection**: Automatically identifies and extracts structured data
- **Confidence Scoring**: Quality assurance for extracted data
- **Fallback Chain**: CSS/XPath extraction when AI fails

### 🏗️ **Microservice Architecture**
- **HTTP API Server**: RESTful endpoints for scraping operations
- **Docker Support**: Container-ready with Docker Compose
- **Kubernetes Ready**: Production deployment manifests included
- **Load Balancing**: Nginx configuration for horizontal scaling

### ⚙️ **Flexible Configuration System**
- **JSON Configuration**: File-based configuration management
- **Environment Variables**: 12-factor app compliance
- **CLI Tools**: Interactive setup and validation
- **Hot Reloading**: Runtime configuration updates

### 🌐 **Multi-Engine Browser Support**
- **ChromeDP**: High-performance Chrome automation
- **Rod**: Lightning-fast browser control
- **Stealth Mode**: Advanced anti-detection techniques
- **Headless & GUI**: Flexible rendering options

### 🚀 **Production Features**
- **Rate Limiting**: Configurable request throttling
- **Caching**: Redis and in-memory caching
- **Proxy Support**: IP rotation and geo-targeting
- **Health Checks**: Monitoring and observability
- **Graceful Shutdown**: Clean resource management

## 📦 Installation

```bash
go get github.com/ramusaaa/goscraper
```

## 🚀 Quick Start

### Method 1: Interactive Setup (Recommended)

```bash
# 1. Initialize configuration
make init-config

# 2. Interactive setup wizard
make setup
# Follow prompts to configure AI keys, caching, etc.

# 3. Validate configuration
make validate-config

# 4. Start the server
make run
```

### Method 2: Environment Variables

```bash
# Set your API keys
export OPENAI_API_KEY="your-openai-key"
export GOSCRAPER_AI_ENABLED=true

# Start the server
go run ./cmd/api
```

### Method 3: Manual Configuration

```bash
# Create config file
cp goscraper.example.json goscraper.json

# Edit configuration
vim goscraper.json

# Start server
go run ./cmd/api
```

## 💻 Usage Examples

### Basic Library Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ramusaaa/goscraper"
)

func main() {
    // Simple scraping
    scraper := goscraper.New()
    
    resp, err := scraper.Get("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    title := resp.Document.Find("title").Text()
    fmt.Printf("Page title: %s\n", title)
}
```

### Advanced Configuration

```go
scraper := goscraper.New(
    goscraper.WithTimeout(30*time.Second),
    goscraper.WithUserAgent("MyBot/1.0"),
    goscraper.WithHeaders(map[string]string{
        "Accept-Language": "en-US,en;q=0.9",
    }),
    goscraper.WithRateLimit(500*time.Millisecond),
    goscraper.WithMaxRetries(3),
    goscraper.WithProxy("http://proxy.example.com:8080"),
    goscraper.WithStealth(true),
)
```

### HTTP API Usage

```bash
# Health check
curl http://localhost:8080/health

# Get configuration
curl http://localhost:8080/config

# Scrape a website
curl -X POST http://localhost:8080/api/scrape \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Smart AI-powered scraping
curl -X POST http://localhost:8080/api/smart-scrape \
  -H "Content-Type: application/json" \
  -d '{"url": "https://shop.example.com/products"}'
```

### Client SDK Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ramusaaa/goscraper/client"
)

func main() {
    // Create client for remote scraper service
    client := client.NewScraperClient("http://localhost:8080")
    
    // Health check
    if err := client.Health(); err != nil {
        log.Fatal("Service unavailable:", err)
    }
    
    // Scrape website
    data, err := client.Scrape("https://example.com")
    if err != nil {
        log.Fatal("Scraping failed:", err)
    }
    
    fmt.Printf("Title: %s\n", data.Title)
    fmt.Printf("Status: %d\n", data.StatusCode)
}
```

## 📋 Configuration Reference

### Configuration File Structure

```json
{
  "server": {
    "port": "8080",
    "host": "0.0.0.0",
    "read_timeout": "30s",
    "write_timeout": "30s"
  },
  "ai": {
    "enabled": true,
    "provider": "openai",
    "confidence_threshold": 0.8,
    "fallback_chain": ["openai", "css", "xpath"],
    "models": {
      "openai": {
        "api_key": "your-openai-key",
        "model": "gpt-4"
      },
      "anthropic": {
        "api_key": "your-anthropic-key",
        "model": "claude-3-sonnet-20240229"
      }
    }
  },
  "browser": {
    "engine": "chromedp",
    "headless": true,
    "stealth": true,
    "pool_size": 5
  },
  "cache": {
    "enabled": true,
    "type": "redis",
    "ttl": "1h",
    "redis": {
      "host": "localhost",
      "port": 6379
    }
  },
  "rate_limit": {
    "requests_per_second": 10,
    "delay": "100ms"
  }
}
```

### Environment Variables

```bash
# Server Configuration
GOSCRAPER_PORT=8080
GOSCRAPER_HOST=0.0.0.0

# AI Configuration
GOSCRAPER_AI_ENABLED=true
GOSCRAPER_AI_PROVIDER=openai
OPENAI_API_KEY=your-openai-key
ANTHROPIC_API_KEY=your-anthropic-key

# Browser Configuration
GOSCRAPER_BROWSER_ENGINE=chromedp
GOSCRAPER_BROWSER_HEADLESS=true
GOSCRAPER_BROWSER_STEALTH=true

# Cache Configuration
GOSCRAPER_CACHE_ENABLED=true
GOSCRAPER_CACHE_TYPE=redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Rate Limiting
GOSCRAPER_RATE_LIMIT_RPS=10
GOSCRAPER_RATE_LIMIT_DELAY=100ms
```

## 🛠️ CLI Tools

### Available Commands

```bash
# Configuration Management
make init-config          # Create default configuration
make setup                # Interactive setup wizard
make validate-config      # Validate configuration
make show-config          # Display current configuration

# Development
make build                # Build binaries
make run                  # Start API server
make test                 # Run tests

# Docker
make docker-build         # Build Docker image
make docker-compose-up    # Start with Docker Compose
make docker-compose-down  # Stop Docker services

# Kubernetes
make k8s-deploy          # Deploy to Kubernetes
make k8s-delete          # Remove from Kubernetes
```

### CLI Usage Examples

```bash
# Initialize new project
goscraper init

# Interactive setup
goscraper setup

# Validate configuration
goscraper validate

# Show current configuration
goscraper config
```
## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │   API Gateway   │    │  Web Dashboard  │
│     (Nginx)     │    │   (Optional)    │    │   (Optional)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Scraper Node 1 │    │  Scraper Node 2 │    │  Scraper Node N │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │HTTP API     │ │    │ │HTTP API     │ │    │ │HTTP API     │ │
│ │Server       │ │    │ │Server       │ │    │ │Server       │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │Browser Pool │ │    │ │Browser Pool │ │    │ │Browser Pool │ │
│ │+ AI Engine  │ │    │ │+ AI Engine  │ │    │ │+ AI Engine  │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
    ┌─────────────────────────────────────────────────────────┐
    │                Infrastructure Layer                      │
    │                                                         │
    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
    │  │    Redis    │  │  Config     │  │   Proxy     │     │
    │  │   Cache     │  │  Storage    │  │  Rotation   │     │
    │  └─────────────┘  └─────────────┘  └─────────────┘     │
    │                                                         │
    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
    │  │   OpenAI    │  │ Anthropic   │  │   Local     │     │
    │  │    API      │  │    API      │  │   Models    │     │
    │  └─────────────┘  └─────────────┘  └─────────────┘     │
    └─────────────────────────────────────────────────────────┘
```

## 🚀 Deployment Options

### 1. Standalone Binary

```bash
# Build and run
go build -o goscraper ./cmd/api
./goscraper
```

### 2. Docker Container

```bash
# Build image
docker build -t goscraper:latest .

# Run container
docker run -p 8080:8080 \
  -e OPENAI_API_KEY=your-key \
  -e GOSCRAPER_AI_ENABLED=true \
  goscraper:latest
```

### 3. Docker Compose

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f scraper-api

# Stop services
docker-compose down
```

### 4. Kubernetes

```bash
# Deploy to cluster
kubectl apply -f k8s/

# Scale deployment
kubectl scale deployment scraper-api --replicas=5

# Check status
kubectl get pods -l app=scraper-api
```

## 🎯 Use Cases & Examples

### E-commerce Price Monitoring

```go
// Monitor product prices
data, err := goscraper.SmartScrape("https://shop.example.com/product/123")
if err != nil {
    log.Fatal(err)
}

if data.ContentType == goscraper.ContentTypeEcommerce {
    for _, product := range data.Products {
        fmt.Printf("Product: %s\n", product.Name)
        fmt.Printf("Price: %s %s\n", product.Price, product.Currency)
        fmt.Printf("Rating: %.1f/5\n", product.Rating)
    }
}
```

### News Aggregation

```go
// Extract news articles
data, err := goscraper.SmartScrape("https://news.example.com/article/123")
if err != nil {
    log.Fatal(err)
}

if data.ContentType == goscraper.ContentTypeNews && data.Article != nil {
    fmt.Printf("Headline: %s\n", data.Article.Headline)
    fmt.Printf("Author: %s\n", data.Article.Author)
    fmt.Printf("Published: %s\n", data.Article.PublishDate)
    fmt.Printf("Content: %s\n", data.Article.Content)
}
```

### Job Listings Scraping

```go
// Extract job postings
data, err := goscraper.SmartScrape("https://jobs.example.com/posting/123")
if err != nil {
    log.Fatal(err)
}

if data.ContentType == goscraper.ContentTypeJob && data.JobListing != nil {
    fmt.Printf("Title: %s\n", data.JobListing.Title)
    fmt.Printf("Company: %s\n", data.JobListing.Company)
    fmt.Printf("Salary: %s\n", data.JobListing.Salary)
    fmt.Printf("Location: %s\n", data.JobListing.Location)
}
```

### Microservice Integration

```go
// Use as microservice client
client := client.NewScraperClient("http://scraper-service:8080")

// Health check
if err := client.Health(); err != nil {
    log.Fatal("Scraper service unavailable")
}

// Batch scraping
urls := []string{
    "https://site1.com",
    "https://site2.com", 
    "https://site3.com",
}

for _, url := range urls {
    data, err := client.Scrape(url)
    if err != nil {
        log.Printf("Failed to scrape %s: %v", url, err)
        continue
    }
    
    fmt.Printf("Scraped %s: %s\n", url, data.Title)
}
```

## 📊 API Reference

### HTTP Endpoints

| Endpoint | Method | Description | Example |
|----------|--------|-------------|---------|
| `/health` | GET | Health check and status | `curl http://localhost:8080/health` |
| `/config` | GET | Current configuration | `curl http://localhost:8080/config` |
| `/api/scrape` | POST | Basic web scraping | See below |
| `/api/smart-scrape` | POST | AI-powered extraction | See below |

### API Examples

#### Basic Scraping

```bash
curl -X POST http://localhost:8080/api/scrape \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "options": {
      "timeout": "30s",
      "user_agent": "Custom Bot"
    }
  }'
```

#### Smart AI Scraping

```bash
curl -X POST http://localhost:8080/api/smart-scrape \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://shop.example.com/product/123"
  }'
```

#### Response Format

```json
{
  "success": true,
  "data": {
    "title": "Page Title",
    "description": "Meta description",
    "url": "https://example.com",
    "status_code": 200,
    "load_time": "1.234s",
    "html": "<!DOCTYPE html>..."
  }
}
```

## 🧪 Testing & Validation

### Run Tests

```bash
# Run all tests
make test

# Run integration tests
go test ./tests/ -v

# Validate configuration
make validate-config

# Full feature validation
./validate_features.sh
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📈 Performance & Monitoring

### Performance Benchmarks

| Feature | Performance | Scalability |
|---------|-------------|-------------|
| **HTTP Requests** | 1,000+ req/sec | Linear scaling |
| **Browser Sessions** | 50+ concurrent | Auto-scaling |
| **AI Extraction** | 10+ pages/sec | Model dependent |
| **Cache Hit Ratio** | 90%+ | Distributed |
| **Memory Usage** | <512MB | Configurable |

### Health Monitoring

```bash
# Check service health
curl http://localhost:8080/health

# Monitor with watch
watch -n 5 'curl -s http://localhost:8080/health | jq'

# Check configuration
curl http://localhost:8080/config | jq
```

### Logging

```bash
# View logs in Docker
docker-compose logs -f scraper-api

# View logs in Kubernetes
kubectl logs -f deployment/scraper-api

# Custom log level
GOSCRAPER_LOG_LEVEL=debug go run ./cmd/api
```

## 🔧 Troubleshooting

### Common Issues

#### Configuration Not Loading

```bash
# Check config file location
go run ./cmd/cli config

# Validate configuration
go run ./cmd/cli validate

# Reset to defaults
go run ./cmd/cli init
```

#### AI Features Not Working

```bash
# Check AI configuration
curl http://localhost:8080/config | jq '.data.ai_enabled'

# Verify API key
export OPENAI_API_KEY=your-key
go run ./cmd/cli validate
```

#### Performance Issues

```bash
# Check resource usage
docker stats scraper-api

# Monitor requests
curl http://localhost:8080/health

# Adjust rate limiting
export GOSCRAPER_RATE_LIMIT_RPS=5
```

### Debug Mode

```bash
# Enable debug logging
export GOSCRAPER_LOG_LEVEL=debug

# Verbose output
go run ./cmd/api -v

# Profile performance
go run ./cmd/api -cpuprofile=cpu.prof
```

## 🤝 Contributing

### Development Setup

```bash
# Clone repository
git clone https://github.com/ramusaaa/goscraper
cd goscraper

# Install dependencies
go mod tidy

# Run tests
make test

# Start development server
make run
```

### Code Style

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Generate documentation
godoc -http=:6060
```

### Pull Request Process

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [ChromeDP](https://github.com/chromedp/chromedp) - Browser automation
- [GoQuery](https://github.com/PuerkitoBio/goquery) - HTML parsing
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP routing
- [Go-Redis](https://github.com/redis/go-redis) - Redis client

## 🏆 Why Choose GoScraper?

| Feature | GoScraper | Scrapy | Puppeteer | Selenium |
|---------|-----------|--------|-----------|----------|
| **Language** | ✅ Go | 🐍 Python | 🟨 JavaScript | 🐍 Python/Java |
| **Performance** | ✅ High | ⚠️ Medium | ⚠️ Medium | ❌ Low |
| **AI Integration** | ✅ Built-in | ❌ External | ❌ External | ❌ External |
| **Microservice Ready** | ✅ Native | ⚠️ Custom | ⚠️ Custom | ❌ No |
| **Configuration** | ✅ Advanced | ⚠️ Basic | ⚠️ Basic | ⚠️ Basic |
| **Stealth Features** | ✅ Advanced | ⚠️ Limited | ✅ Good | ⚠️ Limited |
| **Deployment** | ✅ Easy | ⚠️ Medium | ⚠️ Medium | ❌ Complex |

---

**⭐ Star this repository if you find it useful!**

**🐛 [Report Issues](https://github.com/ramusaaa/goscraper/issues)**

**💡 [Request Features](https://github.com/ramusaaa/goscraper/discussions)**

**💰 [Become a Sponsor](https://github.com/sponsors/ramusaaa)**
