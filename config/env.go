package config

import (
	"os"
	"strconv"
	"time"
)

func (c *Config) LoadFromEnv() {
	
	if port := os.Getenv("GOSCRAPER_PORT"); port != "" {
		c.Server.Port = port
	}
	if host := os.Getenv("GOSCRAPER_HOST"); host != "" {
		c.Server.Host = host
	}

	if enabled := os.Getenv("GOSCRAPER_AI_ENABLED"); enabled != "" {
		c.AI.Enabled = enabled == "true"
	}
	if provider := os.Getenv("GOSCRAPER_AI_PROVIDER"); provider != "" {
		c.AI.Provider = provider
	}

	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		if c.AI.Models == nil {
			c.AI.Models = make(map[string]ModelConfig)
		}
		c.AI.Models["openai"] = ModelConfig{
			APIKey: apiKey,
			Model:  getEnvOrDefault("OPENAI_MODEL", "gpt-4"),
		}
	}

	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		if c.AI.Models == nil {
			c.AI.Models = make(map[string]ModelConfig)
		}
		c.AI.Models["anthropic"] = ModelConfig{
			APIKey: apiKey,
			Model:  getEnvOrDefault("ANTHROPIC_MODEL", "claude-3-sonnet-20240229"),
		}
	}

	if engine := os.Getenv("GOSCRAPER_BROWSER_ENGINE"); engine != "" {
		c.Browser.Engine = engine
	}
	if headless := os.Getenv("GOSCRAPER_BROWSER_HEADLESS"); headless != "" {
		c.Browser.Headless = headless == "true"
	}
	if stealth := os.Getenv("GOSCRAPER_BROWSER_STEALTH"); stealth != "" {
		c.Browser.Stealth = stealth == "true"
	}

	if enabled := os.Getenv("GOSCRAPER_CACHE_ENABLED"); enabled != "" {
		c.Cache.Enabled = enabled == "true"
	}
	if cacheType := os.Getenv("GOSCRAPER_CACHE_TYPE"); cacheType != "" {
		c.Cache.Type = cacheType
	}

	if host := os.Getenv("REDIS_HOST"); host != "" {
		c.Cache.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Cache.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		c.Cache.Redis.Password = password
	}

	if enabled := os.Getenv("GOSCRAPER_PROXY_ENABLED"); enabled != "" {
		c.Proxy.Enabled = enabled == "true"
	}
	if urls := os.Getenv("GOSCRAPER_PROXY_URLS"); urls != "" {

		c.Proxy.URLs = []string{urls} 
	}

	if rps := os.Getenv("GOSCRAPER_RATE_LIMIT_RPS"); rps != "" {
		if r, err := strconv.Atoi(rps); err == nil {
			c.RateLimit.RequestsPerSecond = r
		}
	}
	if delay := os.Getenv("GOSCRAPER_RATE_LIMIT_DELAY"); delay != "" {
		if d, err := time.ParseDuration(delay); err == nil {
			c.RateLimit.Delay = d
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}