package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ramusaaa/routix"
	"github.com/ramusaaa/goscraper"
	"github.com/ramusaaa/goscraper/config"
)

type ScrapeRequest struct {
	URL     string            `json:"url"`
	Options map[string]string `json:"options,omitempty"`
}

type ScrapeResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type APIServer struct {
	scraper *goscraper.GoScraper
	config  *config.Config
}

func NewAPIServer(cfg *config.Config) *APIServer {
	var options []goscraper.Option
	
	options = append(options, goscraper.WithStealth(cfg.Browser.Stealth))
	options = append(options, goscraper.WithTimeout(cfg.Server.ReadTimeout))
	options = append(options, goscraper.WithRateLimit(cfg.RateLimit.Delay))
	
	if cfg.Browser.UserAgent != "" {
		options = append(options, goscraper.WithUserAgent(cfg.Browser.UserAgent))
	}
	
	if cfg.Proxy.Enabled && len(cfg.Proxy.URLs) > 0 {
		options = append(options, goscraper.WithProxy(cfg.Proxy.URLs[0]))
	}

	return &APIServer{
		scraper: goscraper.NewGoScraper(options...),
		config:  cfg,
	}
}

func (s *APIServer) handleScrape(ctx *routix.Context) error {
	var req ScrapeRequest
	if err := ctx.ParseJSON(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ScrapeResponse{
			Success: false,
			Error:   "Invalid JSON",
		})
	}

	if req.URL == "" {
		return ctx.JSON(http.StatusBadRequest, ScrapeResponse{
			Success: false,
			Error:   "URL is required",
		})
	}

	resp, err := s.scraper.Get(req.URL)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ScrapeResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	title := ""
	description := ""
	if resp.Document != nil {
		title = resp.Document.Find("title").Text()
		description, _ = resp.Document.Find("meta[name='description']").Attr("content")
	}

	return ctx.JSON(http.StatusOK, ScrapeResponse{
		Success: true,
		Data: map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         resp.URL,
			"status_code": resp.StatusCode,
			"html":        resp.Body,
			"load_time":   resp.LoadTime.String(),
		},
	})
}

func (s *APIServer) handleSmartScrape(ctx *routix.Context) error {
	var req ScrapeRequest
	if err := ctx.ParseJSON(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ScrapeResponse{
			Success: false,
			Error:   "Invalid JSON",
		})
	}

	data, err := goscraper.SmartScrape(req.URL)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ScrapeResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, ScrapeResponse{
		Success: true,
		Data:    data,
	})
}

func (s *APIServer) handleHealth(ctx *routix.Context) error {
	return ctx.JSON(http.StatusOK, ScrapeResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":     "healthy",
			"time":       time.Now().Format(time.RFC3339),
			"ai_enabled": s.config.AI.Enabled,
			"version":    "1.0.0",
		},
	})
}

func (s *APIServer) handleConfig(ctx *routix.Context) error {
	safeConfig := map[string]interface{}{
		"ai_enabled":     s.config.AI.Enabled,
		"ai_provider":    s.config.AI.Provider,
		"browser_engine": s.config.Browser.Engine,
		"cache_enabled":  s.config.Cache.Enabled,
		"proxy_enabled":  s.config.Proxy.Enabled,
	}
	
	return ctx.JSON(http.StatusOK, ScrapeResponse{
		Success: true,
		Data:    safeConfig,
	})
}



func main() {
	configPath := config.GetConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	fmt.Printf("Loaded config from: %s\n", configPath)
	if cfg.AI.Enabled {
		fmt.Printf("AI enabled with provider: %s\n", cfg.AI.Provider)
	} else {
		fmt.Println("AI disabled - using CSS/XPath extraction only")
	}

	server := NewAPIServer(cfg)
	
	app := routix.New()
	
	corsMiddleware := func(next routix.Handler) routix.Handler {
		return func(ctx *routix.Context) error {
			ctx.SetHeader("Access-Control-Allow-Origin", "*")
			ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type")
			
			if ctx.Request.Method == "OPTIONS" {
				return ctx.String(http.StatusOK, "")
			}
			
			return next(ctx)
		}
	}
	
	app.Use(corsMiddleware)
	
	app.POST("/api/scrape", server.handleScrape)
	app.POST("/api/smart-scrape", server.handleSmartScrape)
	app.GET("/health", server.handleHealth)
	app.GET("/config", server.handleConfig)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Scraper API server starting on %s\n", addr)
	
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      app,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	
	log.Fatal(httpServer.ListenAndServe())
}