package goscraper

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Scraper interface {
	Get(url string) (*Response, error)
	GetWithContext(ctx context.Context, url string) (*Response, error)
	SetConfig(config *Config)
}

type Response struct {
	URL        string
	StatusCode int
	Headers    http.Header
	Body       string
	Document   *goquery.Document
	LoadTime   time.Duration
}

type DefaultScraper struct {
	client *Client
	config *Config
}

func New(options ...Option) *DefaultScraper {
	config := DefaultConfig()
	
	for _, option := range options {
		option(config)
	}

	return &DefaultScraper{
		client: NewClient(config),
		config: config,
	}
}

func (s *DefaultScraper) Get(url string) (*Response, error) {
	return s.GetWithContext(context.Background(), url)
}

func (s *DefaultScraper) GetWithContext(ctx context.Context, url string) (*Response, error) {
	start := time.Now()
	
	resp, err := s.client.GetWithContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	reader := resp.Body
	
	encoding := resp.Header.Get("Content-Encoding")
	if encoding == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	body, _ := doc.Html()
	
	return &Response{
		URL:        url,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
		Document:   doc,
		LoadTime:   time.Since(start),
	}, nil
}

func (s *DefaultScraper) SetConfig(config *Config) {
	s.config = config
	s.client = NewClient(config)
}