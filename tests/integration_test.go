package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/ramusaaa/goscraper/config"
)

func TestConfigSystem(t *testing.T) {
	// Test 1: Default config creation
	cfg := config.DefaultConfig()
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Server.Port)
	}

	// Test 2: Environment variables
	os.Setenv("GOSCRAPER_AI_ENABLED", "true")
	os.Setenv("OPENAI_API_KEY", "test-key")
	cfg.LoadFromEnv()
	
	if !cfg.AI.Enabled {
		t.Error("AI should be enabled from environment variable")
	}

	if len(cfg.AI.Models) == 0 {
		t.Error("OpenAI model should be configured from environment")
	}

	// Test 3: Config validation
	if err := cfg.Validate(); err != nil {
		t.Errorf("Config validation failed: %v", err)
	}
}

func TestAPIEndpoints(t *testing.T) {
	// Note: This test requires the API server to be running
	// Start with: GOSCRAPER_PORT=8084 go run ./cmd/api
	
	baseURL := "http://localhost:8084"
	
	// Test health endpoint
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Logf("API server not running on %s, skipping API tests: %v", baseURL, err)
		t.Skip("API server not running, skipping API tests")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Health check failed with status %d", resp.StatusCode)
	}

	// Test config endpoint
	resp, err = http.Get(baseURL + "/config")
	if err != nil {
		t.Errorf("Config endpoint failed: %v", err)
		return
	}
	defer resp.Body.Close()

	var configResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&configResp); err != nil {
		t.Errorf("Failed to decode config response: %v", err)
	}

	// Test scrape endpoint
	scrapeReq := map[string]string{"url": "https://httpbin.org/html"}
	jsonData, _ := json.Marshal(scrapeReq)

	resp, err = http.Post(
		baseURL+"/api/scrape",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Errorf("Scrape endpoint failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var scrapeResp map[string]interface{}
	if err := json.Unmarshal(body, &scrapeResp); err != nil {
		t.Errorf("Failed to decode scrape response: %v", err)
	}

	if success, ok := scrapeResp["success"].(bool); !ok || !success {
		t.Errorf("Scrape request failed: %v", scrapeResp)
	}
}

func TestCLICommands(t *testing.T) {
	// This would require running CLI commands
	// For now, we'll test the config functions directly
	
	tempFile := "/tmp/test_goscraper_cli.json"
	defer os.Remove(tempFile)

	cfg := config.DefaultConfig()
	cfg.AI.Enabled = true
	
	if err := cfg.Save(tempFile); err != nil {
		t.Errorf("Failed to save config: %v", err)
	}

	loadedCfg, err := config.LoadConfig(tempFile)
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}

	if loadedCfg.AI.Enabled != cfg.AI.Enabled {
		t.Error("Config not loaded correctly")
	}
}