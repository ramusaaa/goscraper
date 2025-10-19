package goscraper

import "time"

func EcommercePreset() []Option {
	return []Option{
		WithStealth(true),
		WithUserAgentRotation(true),
		WithRandomHeaders(true),
		WithHumanDelay(true),
		WithTimeout(45 * time.Second),
		WithRateLimit(3 * time.Second),
		WithMaxRetries(3),
		WithHeaders(map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
			"Accept-Language": "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7",
			"Accept-Encoding": "gzip, deflate, br",
			"DNT":             "1",
			"Connection":      "keep-alive",
			"Sec-Fetch-Dest":  "document",
			"Sec-Fetch-Mode":  "navigate",
			"Sec-Fetch-Site":  "none",
			"Sec-Fetch-User":  "?1",
		}),
	}
}

func NewsPreset() []Option {
	return []Option{
		WithStealth(true),
		WithTimeout(30 * time.Second),
		WithRateLimit(2 * time.Second),
		WithMaxRetries(2),
		WithUserAgent("GoScraper-NewsBot/1.0 (+https://github.com/goscraper/goscraper)"),
	}
}

func SocialMediaPreset() []Option {
	return []Option{
		WithStealth(true),
		WithUserAgentRotation(true),
		WithRandomHeaders(true),
		WithHumanDelay(true),
		WithTimeout(60 * time.Second),
		WithRateLimit(5 * time.Second),
		WithMaxRetries(5),
	}
}

func APIPreset() []Option {
	return []Option{
		WithTimeout(15 * time.Second),
		WithRateLimit(500 * time.Millisecond),
		WithMaxRetries(3),
		WithHeaders(map[string]string{
			"Accept":       "application/json, text/plain, */*",
			"Content-Type": "application/json",
		}),
	}
}

func FastPreset() []Option {
	return []Option{
		WithTimeout(10 * time.Second),
		WithRateLimit(100 * time.Millisecond),
		WithMaxRetries(1),
	}
}

func RobustPreset() []Option {
	return []Option{
		WithStealth(true),
		WithUserAgentRotation(true),
		WithRandomHeaders(true),
		WithHumanDelay(true),
		WithTimeout(120 * time.Second),
		WithRateLimit(10 * time.Second),
		WithMaxRetries(10),
	}
}