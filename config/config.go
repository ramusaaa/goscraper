package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	AI       AIConfig       `json:"ai"`
	Browser  BrowserConfig  `json:"browser"`
	Cache    CacheConfig    `json:"cache,omitempty"`
	Proxy    ProxyConfig    `json:"proxy,omitempty"`
	RateLimit RateLimitConfig `json:"rate_limit"`
}

type ServerConfig struct {
	Port         string        `json:"port"`
	Host         string        `json:"host"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

type AIConfig struct {
	Enabled   bool                    `json:"enabled"`
	Provider  string                  `json:"provider"` // "openai", "anthropic", "local"
	Models    map[string]ModelConfig  `json:"models"`
	Fallback  []string               `json:"fallback_chain"`
	Threshold float64                `json:"confidence_threshold"`
}

type ModelConfig struct {
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	Endpoint string `json:"endpoint,omitempty"`
}

type BrowserConfig struct {
	Engine     string `json:"engine"` // "chromedp", "rod", "playwright"
	Headless   bool   `json:"headless"`
	Stealth    bool   `json:"stealth"`
	UserAgent  string `json:"user_agent,omitempty"`
	PoolSize   int    `json:"pool_size"`
}

type CacheConfig struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"` // "memory", "redis"
	TTL     time.Duration `json:"ttl"`
	Redis   RedisConfig   `json:"redis,omitempty"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type ProxyConfig struct {
	Enabled   bool     `json:"enabled"`
	URLs      []string `json:"urls"`
	Rotation  bool     `json:"rotation"`
}

type RateLimitConfig struct {
	RequestsPerSecond int           `json:"requests_per_second"`
	BurstSize         int           `json:"burst_size"`
	Delay             time.Duration `json:"delay"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "8080",
			Host:         "0.0.0.0",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		AI: AIConfig{
			Enabled:   false, // Disabled by default
			Provider:  "openai",
			Threshold: 0.8,
			Fallback:  []string{"css", "xpath"},
			Models:    make(map[string]ModelConfig),
		},
		Browser: BrowserConfig{
			Engine:   "chromedp",
			Headless: true,
			Stealth:  true,
			PoolSize: 5,
		},
		Cache: CacheConfig{
			Enabled: false,
			Type:    "memory",
			TTL:     1 * time.Hour,
		},
		RateLimit: RateLimitConfig{
			RequestsPerSecond: 10,
			BurstSize:         20,
			Delay:             100 * time.Millisecond,
		},
	}
}

// LoadConfig loads configuration from file or creates default
func LoadConfig(configPath string) (*Config, error) {
	var config *Config
	
	// If config file doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = DefaultConfig()
		if err := config.Save(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		fmt.Printf("Created default config at: %s\n", configPath)
		fmt.Println("Please edit the config file to add your API keys and settings.")
	} else {
		// Load existing config
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		config = &Config{}
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	if config != nil {
		config.LoadFromEnv()
	}

	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save(configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the default config path
func GetConfigPath() string {
	// Try environment variable first
	if path := os.Getenv("GOSCRAPER_CONFIG"); path != "" {
		return path
	}

	// Try current directory
	if _, err := os.Stat("goscraper.json"); err == nil {
		return "goscraper.json"
	}

	// Try home directory
	if home, err := os.UserHomeDir(); err == nil {
		configPath := filepath.Join(home, ".goscraper", "config.json")
		return configPath
	}

	// Fallback to current directory
	return "goscraper.json"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.AI.Enabled {
		if len(c.AI.Models) == 0 {
			return fmt.Errorf("AI is enabled but no models configured")
		}

		for name, model := range c.AI.Models {
			if model.APIKey == "" && model.Endpoint == "" {
				return fmt.Errorf("model %s has no API key or endpoint", name)
			}
		}
	}

	if c.Cache.Enabled && c.Cache.Type == "redis" {
		if c.Cache.Redis.Host == "" {
			return fmt.Errorf("Redis cache enabled but no host specified")
		}
	}

	return nil
}