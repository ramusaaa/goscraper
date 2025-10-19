package goscraper

import (
	"context"
	"time"
)

type GoScraper struct {
	scraper Scraper
}

func NewGoScraper(options ...Option) *GoScraper {
	return &GoScraper{
		scraper: New(options...),
	}
}

func (g *GoScraper) Get(url string) (*Response, error) {
	return g.scraper.Get(url)
}

func (g *GoScraper) GetWithContext(ctx context.Context, url string) (*Response, error) {
	return g.scraper.GetWithContext(ctx, url)
}

func (g *GoScraper) SetConfig(config *Config) {
	g.scraper.SetConfig(config)
}

func QuickScrape(url string) (*Response, error) {
	scraper := New(
		WithStealth(true),
		WithTimeout(30*time.Second),
		WithRateLimit(1*time.Second),
	)
	return scraper.Get(url)
}

func StealthScrape(url string) (*Response, error) {
	scraper := New(
		WithStealth(true),
		WithUserAgentRotation(true),
		WithRandomHeaders(true),
		WithHumanDelay(true),
		WithTimeout(45*time.Second),
		WithRateLimit(2*time.Second),
		WithMaxRetries(3),
	)
	return scraper.Get(url)
}

func SmartScrape(url string) (*SmartData, error) {
	resp, err := StealthScrape(url)
	if err != nil {
		return nil, err
	}
	
	extractor := NewSmartExtractor()
	return extractor.ExtractSmart(resp), nil
}