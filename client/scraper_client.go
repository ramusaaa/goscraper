package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ScraperClient struct {
	baseURL    string
	httpClient *http.Client
}

type ScrapeRequest struct {
	URL     string            `json:"url"`
	Options map[string]string `json:"options,omitempty"`
}

type ScrapeResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ScrapedData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	HTML        string `json:"html"`
}

func NewScraperClient(baseURL string) *ScraperClient {
	return &ScraperClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *ScraperClient) Scrape(url string) (*ScrapedData, error) {
	req := ScrapeRequest{URL: url}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/scrape",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scrapeResp ScrapeResponse
	if err := json.NewDecoder(resp.Body).Decode(&scrapeResp); err != nil {
		return nil, err
	}

	if !scrapeResp.Success {
		return nil, fmt.Errorf("scraping failed: %s", scrapeResp.Error)
	}

	// Convert interface{} to ScrapedData
	dataBytes, err := json.Marshal(scrapeResp.Data)
	if err != nil {
		return nil, err
	}

	var data ScrapedData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *ScraperClient) SmartScrape(url string) (interface{}, error) {
	req := ScrapeRequest{URL: url}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/smart-scrape",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scrapeResp ScrapeResponse
	if err := json.NewDecoder(resp.Body).Decode(&scrapeResp); err != nil {
		return nil, err
	}

	if !scrapeResp.Success {
		return nil, fmt.Errorf("smart scraping failed: %s", scrapeResp.Error)
	}

	return scrapeResp.Data, nil
}

func (c *ScraperClient) Health() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}