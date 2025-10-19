package goscraper

import (
	"net/http"
	"time"
)

type Config struct {
	Timeout         time.Duration
	MaxRedirects    int
	UserAgent       string
	Headers         map[string]string
	Cookies         []*http.Cookie
	
	RateLimit       time.Duration
	MaxConcurrency  int
	
	MaxRetries      int
	RetryDelay      time.Duration
	
	ProxyURL        string
	
	EnableJS        bool
	JSTimeout       time.Duration
	
	EnableStealth   bool
	RotateUA        bool
	RandomHeaders   bool
	HumanDelay      bool
}

type Option func(*Config)

func DefaultConfig() *Config {
	return &Config{
		Timeout:        30 * time.Second,
		MaxRedirects:   10,
		UserAgent:      UserAgentDefault,
		Headers:        make(map[string]string),
		RateLimit:      100 * time.Millisecond,
		MaxConcurrency: 10,
		MaxRetries:     3,
		RetryDelay:     1 * time.Second,
		EnableJS:       false,
		JSTimeout:      10 * time.Second,
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Config) {
		c.UserAgent = userAgent
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *Config) {
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

func WithRateLimit(delay time.Duration) Option {
	return func(c *Config) {
		c.RateLimit = delay
	}
}

func WithMaxRetries(retries int) Option {
	return func(c *Config) {
		c.MaxRetries = retries
	}
}

func WithProxy(proxyURL string) Option {
	return func(c *Config) {
		c.ProxyURL = proxyURL
	}
}

func WithJavaScript(enabled bool) Option {
	return func(c *Config) {
		c.EnableJS = enabled
	}
}

func WithStealth(enabled bool) Option {
	return func(c *Config) {
		c.EnableStealth = enabled
	}
}

func WithUserAgentRotation(enabled bool) Option {
	return func(c *Config) {
		c.RotateUA = enabled
	}
}

func WithRandomHeaders(enabled bool) Option {
	return func(c *Config) {
		c.RandomHeaders = enabled
	}
}

func WithHumanDelay(enabled bool) Option {
	return func(c *Config) {
		c.HumanDelay = enabled
	}
}