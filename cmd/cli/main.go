package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ramusaaa/goscraper/config"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		initConfig()
	case "config":
		showConfig()
	case "setup":
		setupWizard()
	case "validate":
		validateConfig()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("GoScraper CLI")
	fmt.Println("Usage:")
	fmt.Println("  goscraper init     - Create default config file")
	fmt.Println("  goscraper config   - Show current config")
	fmt.Println("  goscraper setup    - Interactive setup wizard")
	fmt.Println("  goscraper validate - Validate config file")
}

func initConfig() {
	configPath := config.GetConfigPath()
	
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at: %s\n", configPath)
		fmt.Print("Overwrite? (y/N): ")
		
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return
		}
	}

	cfg := config.DefaultConfig()
	if err := cfg.Save(configPath); err != nil {
		fmt.Printf("Error creating config: %v\n", err)
		return
	}

	fmt.Printf("Created config file at: %s\n", configPath)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Edit the config file to add your API keys")
	fmt.Println("2. Run 'goscraper validate' to check your config")
	fmt.Println("3. Run 'goscraper setup' for interactive configuration")
}

func showConfig() {
	configPath := config.GetConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	fmt.Printf("Config loaded from: %s\n\n", configPath)
	fmt.Printf("Server: %s:%s\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("AI Enabled: %v\n", cfg.AI.Enabled)
	if cfg.AI.Enabled {
		fmt.Printf("AI Provider: %s\n", cfg.AI.Provider)
		fmt.Printf("Models configured: %d\n", len(cfg.AI.Models))
	}
	fmt.Printf("Browser Engine: %s\n", cfg.Browser.Engine)
	fmt.Printf("Cache Enabled: %v\n", cfg.Cache.Enabled)
	fmt.Printf("Proxy Enabled: %v\n", cfg.Proxy.Enabled)
}

func setupWizard() {
	fmt.Println("GoScraper Setup Wizard")
	fmt.Println("======================")

	reader := bufio.NewReader(os.Stdin)
	cfg := config.DefaultConfig()

	fmt.Print("\nEnable AI-powered extraction? (y/N): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	cfg.AI.Enabled = response == "y" || response == "yes"

	if cfg.AI.Enabled {
		fmt.Print("AI Provider (openai/anthropic/local): ")
		provider, _ := reader.ReadString('\n')
		cfg.AI.Provider = strings.TrimSpace(provider)

		if cfg.AI.Provider == "openai" || cfg.AI.Provider == "" {
			cfg.AI.Provider = "openai"
			fmt.Print("OpenAI API Key: ")
			apiKey, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)
			
			if apiKey != "" {
				cfg.AI.Models["openai"] = config.ModelConfig{
					APIKey: apiKey,
					Model:  "gpt-4",
				}
			}
		}

		if cfg.AI.Provider == "anthropic" {
			fmt.Print("Anthropic API Key: ")
			apiKey, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)
			
			if apiKey != "" {
				cfg.AI.Models["anthropic"] = config.ModelConfig{
					APIKey: apiKey,
					Model:  "claude-3-sonnet-20240229",
				}
			}
		}
	}

	fmt.Print("\nEnable caching? (y/N): ")
	response, _ = reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	cfg.Cache.Enabled = response == "y" || response == "yes"

	if cfg.Cache.Enabled {
		fmt.Print("Cache type (memory/redis): ")
		cacheType, _ := reader.ReadString('\n')
		cfg.Cache.Type = strings.TrimSpace(cacheType)

		if cfg.Cache.Type == "redis" {
			fmt.Print("Redis host (localhost): ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			if host == "" {
				host = "localhost"
			}
			cfg.Cache.Redis.Host = host

			fmt.Print("Redis port (6379): ")
			port, _ := reader.ReadString('\n')
			port = strings.TrimSpace(port)
			if port == "" {
				cfg.Cache.Redis.Port = 6379
			}
		}
	}

	configPath := config.GetConfigPath()
	if err := cfg.Save(configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	fmt.Printf("\nConfiguration saved to: %s\n", configPath)
	fmt.Println("Setup complete! You can now start the scraper.")
}

func validateConfig() {
	configPath := config.GetConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if err := cfg.Validate(); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
		return
	}

	fmt.Println("Configuration is valid! ✓")
	
	if cfg.AI.Enabled {
		fmt.Printf("✓ AI enabled with %d model(s)\n", len(cfg.AI.Models))
	}
	
	if cfg.Cache.Enabled {
		fmt.Printf("✓ Cache enabled (%s)\n", cfg.Cache.Type)
	}
	
	if cfg.Proxy.Enabled {
		fmt.Printf("✓ Proxy enabled with %d URL(s)\n", len(cfg.Proxy.URLs))
	}
}