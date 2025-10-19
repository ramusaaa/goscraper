package main

import (
	"fmt"
	"log"

	"github.com/ramusaaa/goscraper/client"
)

func main() {
	// Scraper microservice client'ı oluştur
	scraperClient := client.NewScraperClient("http://localhost:8080")

	// Health check
	if err := scraperClient.Health(); err != nil {
		log.Fatal("Scraper service is not healthy:", err)
	}

	// Basit scraping
	data, err := scraperClient.Scrape("https://example.com")
	if err != nil {
		log.Fatal("Scraping failed:", err)
	}

	fmt.Printf("Title: %s\n", data.Title)
	fmt.Printf("Description: %s\n", data.Description)

	// Smart scraping (AI-powered)
	smartData, err := scraperClient.SmartScrape("https://trendyol.com/product-url")
	if err != nil {
		log.Fatal("Smart scraping failed:", err)
	}

	fmt.Printf("Smart Data: %+v\n", smartData)
}

// Web framework ile kullanım (Gin örneği)
func webFrameworkExample() {
	// gin := gin.Default()
	// scraperClient := client.NewScraperClient("http://scraper-service:8080")
	
	// gin.POST("/scrape", func(c *gin.Context) {
	//     var req struct {
	//         URL string `json:"url"`
	//     }
	//     
	//     if err := c.ShouldBindJSON(&req); err != nil {
	//         c.JSON(400, gin.H{"error": err.Error()})
	//         return
	//     }
	//     
	//     data, err := scraperClient.Scrape(req.URL)
	//     if err != nil {
	//         c.JSON(500, gin.H{"error": err.Error()})
	//         return
	//     }
	//     
	//     c.JSON(200, data)
	// })
}