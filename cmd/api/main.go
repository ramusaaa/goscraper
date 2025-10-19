package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func (s *APIServer) handleScrape(w http.ResponseWriter, r *http.Request) {
	var req ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		s.sendError(w, "URL is required", http.StatusBadRequest)
		return
	}

	resp, err := s.scraper.Get(req.URL)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := ""
	description := ""
	if resp.Document != nil {
		title = resp.Document.Find("title").Text()
		description, _ = resp.Document.Find("meta[name='description']").Attr("content")
	}

	s.sendSuccess(w, map[string]interface{}{
		"title":       title,
		"description": description,
		"url":         resp.URL,
		"status_code": resp.StatusCode,
		"html":        resp.Body,
		"load_time":   resp.LoadTime.String(),
	})
}

func (s *APIServer) handleSmartScrape(w http.ResponseWriter, r *http.Request) {
	var req ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	data, err := goscraper.SmartScrape(req.URL)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendSuccess(w, data)
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.sendSuccess(w, map[string]interface{}{
		"status":     "healthy",
		"time":       time.Now().Format(time.RFC3339),
		"ai_enabled": s.config.AI.Enabled,
		"version":    "1.0.0",
	})
}

func (s *APIServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	safeConfig := map[string]interface{}{
		"ai_enabled":    s.config.AI.Enabled,
		"ai_provider":   s.config.AI.Provider,
		"browser_engine": s.config.Browser.Engine,
		"cache_enabled": s.config.Cache.Enabled,
		"proxy_enabled": s.config.Proxy.Enabled,
	}
	
	s.sendSuccess(w, safeConfig)
}

func (s *APIServer) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ScrapeResponse{
		Success: true,
		Data:    data,
	})
}

func (s *APIServer) sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ScrapeResponse{
		Success: false,
		Error:   message,
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
	
	r := mux.NewRouter()
	
	r.HandleFunc("/api/scrape", server.handleScrape).Methods("POST")
	r.HandleFunc("/api/smart-scrape", server.handleSmartScrape).Methods("POST")
	r.HandleFunc("/health", server.handleHealth).Methods("GET")
	r.HandleFunc("/config", server.handleConfig).Methods("GET")
	
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			if r.Method == "OPTIONS" {
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Scraper API server starting on %s\n", addr)
	
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	
	log.Fatal(httpServer.ListenAndServe())
}