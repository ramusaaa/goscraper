package internal

import (
	"net/url"
	"strings"
)

func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func NormalizeURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	
	return u.String()
}

func CleanText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}
	
	return text
}

func ExtractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Host
}