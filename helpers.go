package goscraper

import (
	"regexp"
	"strings"
)



type ExtractedData struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Links       []Link            `json:"links"`
	Images      []Image           `json:"images"`
	MetaTags    map[string]string `json:"meta_tags"`
	Text        []string          `json:"text"`
	Emails      []string          `json:"emails"`
	PhoneNumbers []string         `json:"phone_numbers"`
}

func ExtractAll(resp *Response) *ExtractedData {
	parser := NewParser(resp.Document)
	
	return &ExtractedData{
		Title:       parser.ExtractTitle(),
		Description: getMetaDescription(parser),
		Links:       parser.ExtractLinks(),
		Images:      parser.ExtractImages(),
		MetaTags:    parser.ExtractMetaTags(),
		Text:        extractMeaningfulText(parser),
		Emails:      extractEmails(resp.Body),
		PhoneNumbers: extractPhoneNumbers(resp.Body),
	}
}

func ExtractProducts(resp *Response, selectors ProductSelectors) []Product {
	parser := NewParser(resp.Document)
	
	names := parser.ExtractTexts(selectors.Name)
	prices := parser.ExtractTexts(selectors.Price)
	images := parser.ExtractAttrs(selectors.Image, "src")
	links := parser.ExtractAttrs(selectors.Link, "href")
	
	maxLen := max(max(len(names), len(prices)), max(len(images), len(links)))
	products := make([]Product, maxLen)
	
	for i := 0; i < maxLen; i++ {
		product := Product{}
		
		if i < len(names) {
			product.Name = strings.TrimSpace(names[i])
		}
		if i < len(prices) {
			product.Price = strings.TrimSpace(prices[i])
		}
		if i < len(images) {
			product.ImageURL = images[i]
		}
		if i < len(links) {
			product.URL = links[i]
		}
		
		products[i] = product
	}
	
	return products
}

type ProductSelectors struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

type Product struct {
	Name     string `json:"name"`
	Price    string `json:"price"`
	ImageURL string `json:"image_url"`
	URL      string `json:"url"`
}

func GetTrendyolSelectors() ProductSelectors {
	return ProductSelectors{
		Name:  ".prdct-desc-cntnr-name, .product-down .name",
		Price: ".price-current, .prc-box-dscntd",
		Image: ".p-card-img img",
		Link:  ".p-card-wrppr a",
	}
}

func GetHepsiburadaSelectors() ProductSelectors {
	return ProductSelectors{
		Name:  ".product-title, [data-test-id='product-card-name']",
		Price: ".price-current, .currentPrice",
		Image: ".product-image img",
		Link:  ".product-item a",
	}
}

func GetN11Selectors() ProductSelectors {
	return ProductSelectors{
		Name:  ".productName, .pro .productTitle",
		Price: ".newPrice, .priceContainer .newPrice",
		Image: ".productImage img",
		Link:  ".pro a",
	}
}

func getMetaDescription(parser *Parser) string {
	meta := parser.ExtractMetaTags()
	if desc, exists := meta["description"]; exists {
		return desc
	}
	return ""
}

func extractMeaningfulText(parser *Parser) []string {
	selectors := []string{"h1", "h2", "h3", "p", ".title", ".description"}
	var texts []string
	
	for _, selector := range selectors {
		elements := parser.ExtractTexts(selector)
		for _, text := range elements {
			cleaned := strings.TrimSpace(text)
			if len(cleaned) > 10 && len(cleaned) < 500 {
				texts = append(texts, cleaned)
			}
		}
	}
	
	return texts
}

func extractEmails(html string) []string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := emailRegex.FindAllString(html, -1)
	
	unique := make(map[string]bool)
	var emails []string
	
	for _, email := range matches {
		if !unique[email] {
			unique[email] = true
			emails = append(emails, email)
		}
	}
	
	return emails
}

func extractPhoneNumbers(html string) []string {
	phoneRegex := regexp.MustCompile(`(\+90|0)?\s?[0-9]{3}\s?[0-9]{3}\s?[0-9]{2}\s?[0-9]{2}`)
	matches := phoneRegex.FindAllString(html, -1)
	
	unique := make(map[string]bool)
	var phones []string
	
	for _, phone := range matches {
		cleaned := strings.ReplaceAll(phone, " ", "")
		if !unique[cleaned] && len(cleaned) >= 10 {
			unique[cleaned] = true
			phones = append(phones, cleaned)
		}
	}
	
	return phones
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}