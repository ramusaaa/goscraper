/*
Package goscraper provides a modern, fast, and stealth web scraping library for Go.

GoScraper is designed for production use with advanced anti-bot detection evasion,
automatic content extraction, and built-in support for popular e-commerce sites.

# Quick Start

The simplest way to get started is with QuickScrape:

	resp, err := goscraper.QuickScrape("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	
	data := goscraper.ExtractAll(resp)
	fmt.Printf("Title: %s\n", data.Title)

# Stealth Mode

For sites with bot protection, use stealth mode:

	resp, err := goscraper.StealthScrape("https://protected-site.com")
	if err != nil {
		log.Fatal(err)
	}

# E-commerce Scraping

GoScraper includes built-in support for popular e-commerce sites:

	scraper := goscraper.New(goscraper.EcommercePreset()...)
	resp, err := scraper.Get("https://shop.example.com")
	
	products := goscraper.ExtractProducts(resp, goscraper.GetTrendyolSelectors())
	for _, product := range products {
		fmt.Printf("%s - %s\n", product.Name, product.Price)
	}

# Advanced Configuration

For full control, use the configuration options:

	scraper := goscraper.New(
		goscraper.WithStealth(true),
		goscraper.WithUserAgentRotation(true),
		goscraper.WithRandomHeaders(true),
		goscraper.WithHumanDelay(true),
		goscraper.WithTimeout(30*time.Second),
		goscraper.WithRateLimit(2*time.Second),
		goscraper.WithMaxRetries(3),
	)

# Presets

GoScraper provides several presets for common use cases:

	- EcommercePreset(): Optimized for e-commerce sites
	- NewsPreset(): Optimized for news websites
	- SocialMediaPreset(): Optimized for social media platforms
	- APIPreset(): Optimized for API endpoints
	- FastPreset(): Optimized for speed
	- RobustPreset(): Optimized for reliability

# Parsing

The library includes powerful parsing utilities:

	parser := goscraper.NewParser(resp.Document)
	
	title := parser.ExtractTitle()
	links := parser.ExtractLinks()
	images := parser.ExtractImages()
	meta := parser.ExtractMetaTags()

# Error Handling

GoScraper provides detailed error information and automatic retry mechanisms:

	scraper := goscraper.New(
		goscraper.WithMaxRetries(3),
		goscraper.WithRetryDelay(time.Second),
	)

For more examples and documentation, visit: https://github.com/ramusaaa/goscraper
*/
package goscraper