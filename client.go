package goscraper

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ramusaaa/goscraper/pkg/stealth"
)

type Client struct {
	httpClient    *http.Client
	config        *Config
	lastReq       time.Time
	stealthClient *stealth.BotDetectionEvasion
}

func NewClient(config *Config) *Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.MaxRedirects {
				return fmt.Errorf("stopped after %d redirects", config.MaxRedirects)
			}
			return nil
		},
	}

	return &Client{
		httpClient:    client,
		config:        config,
		stealthClient: stealth.NewBotDetectionEvasion(),
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.GetWithContext(context.Background(), url)
}

func (c *Client) GetWithContext(ctx context.Context, url string) (*http.Response, error) {
	c.applyRateLimit()

	if c.config.EnableStealth {
		return c.stealthClient.MakeRequest(url)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range c.config.Cookies {
		req.AddCookie(cookie)
	}

	var resp *http.Response
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		if attempt < c.config.MaxRetries {
			time.Sleep(c.config.RetryDelay * time.Duration(attempt+1))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.MaxRetries+1, err)
	}

	return resp, nil
}

func (c *Client) applyRateLimit() {
	if c.config.RateLimit > 0 {
		elapsed := time.Since(c.lastReq)
		if elapsed < c.config.RateLimit {
			time.Sleep(c.config.RateLimit - elapsed)
		}
		c.lastReq = time.Now()
	}
}