package goscraper

import (
	"regexp"
	"strings"
)

func (se *SmartExtractor) extractProducts(parser *Parser, url string) []SmartProduct {
	var products []SmartProduct
	
	domain := extractDomainFromURL(url)
	if selectors := getProductSelectorsForDomain(domain); selectors != nil {
		return se.extractProductsWithSelectors(parser, *selectors)
	}
	
	productSelectors := []string{
		"[itemtype*='Product']",
		"[data-testid*='product']",
		"[class*='product']",
		"[class*='item']",
		".product, .item, .listing",
	}
	
	for _, selector := range productSelectors {
		elements := parser.ExtractTexts(selector)
		if len(elements) > 0 {
			for i, element := range elements {
				if i >= 20 {
					break
				}
				
				product := SmartProduct{
					Name:    cleanText(element),
					InStock: true,
				}
				
				products = append(products, product)
			}
			break
		}
	}
	
	return products
}

func (se *SmartExtractor) extractArticle(parser *Parser) *Article {
	article := &Article{}
	
	headlines := []string{"h1", ".headline", ".title", "[class*='headline']"}
	for _, selector := range headlines {
		if title := parser.ExtractText(selector); title != "" {
			article.Headline = cleanText(title)
			break
		}
	}
	
	authors := []string{".author", ".byline", "[class*='author']", "[rel='author']"}
	for _, selector := range authors {
		if author := parser.ExtractText(selector); author != "" {
			article.Author = cleanText(author)
			break
		}
	}
	
	contents := []string{".content", ".article-body", ".post-content", "article", ".entry-content"}
	for _, selector := range contents {
		if content := parser.ExtractText(selector); content != "" {
			article.Content = cleanText(content)
			break
		}
	}
	
	dates := []string{".date", ".publish-date", "[datetime]", "time"}
	for _, selector := range dates {
		if date := parser.ExtractText(selector); date != "" {
			article.PublishDate = cleanText(date)
			break
		}
	}
	
	return article
}

func (se *SmartExtractor) extractBlogPost(parser *Parser) *BlogPost {
	post := &BlogPost{}
	
	if title := parser.ExtractTitle(); title != "" {
		post.Title = title
	}
	
	authors := []string{".author", ".post-author", "[class*='author']"}
	for _, selector := range authors {
		if author := parser.ExtractText(selector); author != "" {
			post.Author = cleanText(author)
			break
		}
	}
	
	contents := []string{".post-content", ".entry-content", ".blog-content", "article"}
	for _, selector := range contents {
		if content := parser.ExtractText(selector); content != "" {
			post.Content = cleanText(content)
			break
		}
	}
	
	categories := parser.ExtractTexts(".category, .categories, [class*='category']")
	post.Categories = cleanTextArray(categories)
	
	tags := parser.ExtractTexts(".tag, .tags, [class*='tag']")
	post.Tags = cleanTextArray(tags)
	
	return post
}

func (se *SmartExtractor) extractJobListing(parser *Parser) *JobListing {
	job := &JobListing{}
	
	titles := []string{"h1", ".job-title", ".position", "[class*='title']"}
	for _, selector := range titles {
		if title := parser.ExtractText(selector); title != "" {
			job.Title = cleanText(title)
			break
		}
	}
	
	companies := []string{".company", ".employer", "[class*='company']"}
	for _, selector := range companies {
		if company := parser.ExtractText(selector); company != "" {
			job.Company = cleanText(company)
			break
		}
	}
	
	locations := []string{".location", ".city", "[class*='location']"}
	for _, selector := range locations {
		if location := parser.ExtractText(selector); location != "" {
			job.Location = cleanText(location)
			break
		}
	}
	
	salaries := []string{".salary", ".wage", "[class*='salary']", "[class*='wage']"}
	for _, selector := range salaries {
		if salary := parser.ExtractText(selector); salary != "" {
			job.Salary = cleanText(salary)
			break
		}
	}
	
	descriptions := []string{".description", ".job-description", ".details"}
	for _, selector := range descriptions {
		if desc := parser.ExtractText(selector); desc != "" {
			job.Description = cleanText(desc)
			break
		}
	}
	
	return job
}

func (se *SmartExtractor) extractProperty(parser *Parser) *Property {
	property := &Property{}
	
	if title := parser.ExtractTitle(); title != "" {
		property.Title = title
	}
	
	prices := []string{".price", ".cost", "[class*='price']"}
	for _, selector := range prices {
		if price := parser.ExtractText(selector); price != "" {
			property.Price = cleanText(price)
			break
		}
	}
	
	locations := []string{".location", ".address", "[class*='location']"}
	for _, selector := range locations {
		if location := parser.ExtractText(selector); location != "" {
			property.Location = cleanText(location)
			break
		}
	}
	
	bedrooms := []string{".bedroom", ".bed", "[class*='bedroom']"}
	for _, selector := range bedrooms {
		if bed := parser.ExtractText(selector); bed != "" {
			property.Bedrooms = cleanText(bed)
			break
		}
	}
	
	return property
}

func (se *SmartExtractor) extractRecipe(parser *Parser) *Recipe {
	recipe := &Recipe{}
	
	if name := parser.ExtractTitle(); name != "" {
		recipe.Name = name
	}
	
	ingredients := parser.ExtractTexts(".ingredient, .ingredients li, [class*='ingredient']")
	recipe.Ingredients = cleanTextArray(ingredients)
	
	instructions := parser.ExtractTexts(".instruction, .instructions li, .step, [class*='instruction']")
	recipe.Instructions = cleanTextArray(instructions)
	
	times := []string{".prep-time", ".cook-time", ".total-time", "[class*='time']"}
	for _, selector := range times {
		if time := parser.ExtractText(selector); time != "" {
			if strings.Contains(selector, "prep") {
				recipe.PrepTime = cleanText(time)
			} else if strings.Contains(selector, "cook") {
				recipe.CookTime = cleanText(time)
			} else if strings.Contains(selector, "total") {
				recipe.TotalTime = cleanText(time)
			}
		}
	}
	
	return recipe
}

func (se *SmartExtractor) extractEvent(parser *Parser) *Event {
	event := &Event{}
	
	if name := parser.ExtractTitle(); name != "" {
		event.Name = name
	}
	
	dates := []string{".date", ".event-date", "[class*='date']", "time"}
	for _, selector := range dates {
		if date := parser.ExtractText(selector); date != "" {
			event.Date = cleanText(date)
			break
		}
	}
	
	venues := []string{".venue", ".location", "[class*='venue']"}
	for _, selector := range venues {
		if venue := parser.ExtractText(selector); venue != "" {
			event.Venue = cleanText(venue)
			break
		}
	}
	
	prices := []string{".price", ".ticket-price", "[class*='price']"}
	for _, selector := range prices {
		if price := parser.ExtractText(selector); price != "" {
			event.Price = cleanText(price)
			break
		}
	}
	
	return event
}

func (se *SmartExtractor) extractVideo(parser *Parser) *Video {
	video := &Video{}
	
	if title := parser.ExtractTitle(); title != "" {
		video.Title = title
	}
	
	durations := []string{".duration", "[class*='duration']", ".time"}
	for _, selector := range durations {
		if duration := parser.ExtractText(selector); duration != "" {
			video.Duration = cleanText(duration)
			break
		}
	}
	
	views := []string{".views", "[class*='view']", ".watch-count"}
	for _, selector := range views {
		if view := parser.ExtractText(selector); view != "" {
			video.Views = cleanText(view)
			break
		}
	}
	
	authors := []string{".channel", ".author", "[class*='channel']"}
	for _, selector := range authors {
		if author := parser.ExtractText(selector); author != "" {
			video.Author = cleanText(author)
			break
		}
	}
	
	return video
}

func getProductSelectorsForDomain(domain string) *ProductSelectors {
	domain = strings.ToLower(domain)
	
	if strings.Contains(domain, "trendyol") {
		return &ProductSelectors{
			Name:  ".prdct-desc-cntnr-name, .product-down .name",
			Price: ".price-current, .prc-box-dscntd",
			Image: ".p-card-img img",
			Link:  ".p-card-wrppr a",
		}
	}
	
	if strings.Contains(domain, "hepsiburada") {
		return &ProductSelectors{
			Name:  ".product-title, [data-test-id='product-card-name']",
			Price: ".price-current, .currentPrice",
			Image: ".product-image img",
			Link:  ".product-item a",
		}
	}
	
	if strings.Contains(domain, "n11") {
		return &ProductSelectors{
			Name:  ".productName, .pro .productTitle",
			Price: ".newPrice, .priceContainer .newPrice",
			Image: ".productImage img",
			Link:  ".pro a",
		}
	}
	
	if strings.Contains(domain, "amazon") {
		return &ProductSelectors{
			Name:  "[data-cy='title-recipe-card'], .s-title-instructions-style",
			Price: ".a-price-whole, .a-offscreen",
			Image: ".s-image",
			Link:  ".s-link-style a",
		}
	}
	
	if strings.Contains(domain, "ebay") {
		return &ProductSelectors{
			Name:  ".s-item__title",
			Price: ".s-item__price",
			Image: ".s-item__image",
			Link:  ".s-item__link",
		}
	}
	
	return nil
}

func (se *SmartExtractor) extractProductsWithSelectors(parser *Parser, selectors ProductSelectors) []SmartProduct {
	names := parser.ExtractTexts(selectors.Name)
	prices := parser.ExtractTexts(selectors.Price)
	images := parser.ExtractAttrs(selectors.Image, "src")
	links := parser.ExtractAttrs(selectors.Link, "href")
	
	maxLen := max(max(len(names), len(prices)), max(len(images), len(links)))
	products := make([]SmartProduct, 0, maxLen)
	
	for i := 0; i < maxLen; i++ {
		product := SmartProduct{InStock: true}
		
		if i < len(names) {
			product.Name = cleanText(names[i])
		}
		if i < len(prices) {
			product.Price = extractPrice(prices[i])
			product.Currency = extractCurrency(prices[i])
		}
		if i < len(images) {
			product.ImageURL = images[i]
		}
		if i < len(links) {
			product.URL = links[i]
		}
		
		if product.Name != "" {
			products = append(products, product)
		}
	}
	
	return products
}

func cleanText(text string) string {
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return text
}

func cleanTextArray(texts []string) []string {
	var cleaned []string
	for _, text := range texts {
		if clean := cleanText(text); clean != "" && len(clean) > 2 {
			cleaned = append(cleaned, clean)
		}
	}
	return cleaned
}

func extractPrice(text string) string {
	patterns := []string{
		`\d+[.,]\d+\s*(?:TL|₺|USD|\$|EUR|€)`,
		`(?:TL|₺|USD|\$|EUR|€)\s*\d+[.,]\d+`,
		`\d+[.,]\d+`,
		`\d+`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindString(text); match != "" {
			return strings.TrimSpace(match)
		}
	}
	
	return text
}

func extractCurrency(text string) string {
	currencies := map[string]string{
		"TL": "TRY", "₺": "TRY",
		"$": "USD", "USD": "USD",
		"€": "EUR", "EUR": "EUR",
		"£": "GBP", "GBP": "GBP",
	}
	
	for symbol, code := range currencies {
		if strings.Contains(text, symbol) {
			return code
		}
	}
	
	return "TRY"
}