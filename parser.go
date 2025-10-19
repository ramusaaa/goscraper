package goscraper

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Parser struct {
	doc *goquery.Document
}

func NewParser(doc *goquery.Document) *Parser {
	return &Parser{doc: doc}
}

func (p *Parser) ExtractText(selector string) string {
	return strings.TrimSpace(p.doc.Find(selector).First().Text())
}

func (p *Parser) ExtractTexts(selector string) []string {
	var texts []string
	p.doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			texts = append(texts, text)
		}
	})
	return texts
}

func (p *Parser) ExtractAttr(selector, attr string) string {
	val, _ := p.doc.Find(selector).First().Attr(attr)
	return val
}

func (p *Parser) ExtractAttrs(selector, attr string) []string {
	var attrs []string
	p.doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if val, exists := s.Attr(attr); exists {
			attrs = append(attrs, val)
		}
	})
	return attrs
}

func (p *Parser) ExtractLinks() []Link {
	var links []Link
	p.doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		text := strings.TrimSpace(s.Text())
		links = append(links, Link{
			URL:  href,
			Text: text,
		})
	})
	return links
}

func (p *Parser) ExtractImages() []Image {
	var images []Image
	p.doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		alt, _ := s.Attr("alt")
		images = append(images, Image{
			URL: src,
			Alt: alt,
		})
	})
	return images
}

func (p *Parser) ExtractMetaTags() map[string]string {
	meta := make(map[string]string)
	
	p.doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			if content, exists := s.Attr("content"); exists {
				meta[name] = content
			}
		}
		if property, exists := s.Attr("property"); exists {
			if content, exists := s.Attr("content"); exists {
				meta[property] = content
			}
		}
	})
	
	return meta
}

func (p *Parser) ExtractTitle() string {
	return p.ExtractText("title")
}

func (p *Parser) ExtractByRegex(pattern string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	
	html, _ := p.doc.Html()
	return re.FindAllString(html, -1)
}

type Link struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}