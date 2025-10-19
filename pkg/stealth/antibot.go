package stealth

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type StealthConfig struct {
	RotateUserAgents    bool
	RandomizeHeaders    bool
	SimulateHuman       bool
	UseProxyRotation    bool
	BypassCloudflare    bool
	DelayRange          [2]int
	MaxRetries          int
	TLSFingerprinting   bool
	JSChallengeBypass   bool
}

type StealthClient struct {
	config     *StealthConfig
	userAgents []string
	proxies    []string
	client     *http.Client
}

func NewStealthClient(config *StealthConfig) *StealthClient {
	return &StealthClient{
		config:     config,
		userAgents: getRealisticUserAgents(),
		client:     createStealthHTTPClient(config),
	}
}

func createStealthHTTPClient(config *StealthConfig) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   45 * time.Second,
	}
}

func (s *StealthClient) CreateStealthRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if s.config.RotateUserAgents {
		req.Header.Set("User-Agent", s.getRandomUserAgent())
	}

	if s.config.RandomizeHeaders {
		s.addRealisticHeaders(req)
	}

	return req, nil
}

func (s *StealthClient) getRandomUserAgent() string {
	return s.userAgents[rand.Intn(len(s.userAgents))]
}

func (s *StealthClient) addRealisticHeaders(req *http.Request) {
	headers := map[string][]string{
		"Accept": {
			"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		},
		"Accept-Language": {
			"tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7",
			"tr,en-US;q=0.9,en;q=0.8",
			"tr-TR,tr;q=0.8,en-US;q=0.5,en;q=0.3",
		},
		"Accept-Encoding": {
			"gzip, deflate, br",
			"gzip, deflate",
		},
		"Cache-Control": {
			"max-age=0",
			"no-cache",
			"",
		},
		"Sec-Fetch-Dest": {
			"document",
			"empty",
		},
		"Sec-Fetch-Mode": {
			"navigate",
			"cors",
		},
		"Sec-Fetch-Site": {
			"none",
			"same-origin",
			"cross-site",
		},
		"Sec-Fetch-User": {
			"?1",
			"",
		},
	}

	for header, options := range headers {
		if len(options) > 0 {
			value := options[rand.Intn(len(options))]
			if value != "" {
				req.Header.Set(header, value)
			}
		}
	}

	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	if rand.Float32() < 0.3 {
		req.Header.Set("Sec-CH-UA", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		req.Header.Set("Sec-CH-UA-Mobile", "?0")
		req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	}
}

func (s *StealthClient) SimulateHumanDelay() {
	if s.config.SimulateHuman {
		min := s.config.DelayRange[0]
		max := s.config.DelayRange[1]
		delay := time.Duration(min+rand.Intn(max-min)) * time.Millisecond
		time.Sleep(delay)
	}
}

func getRealisticUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Android 14; Mobile; rv:121.0) Gecko/121.0 Firefox/121.0",
		"Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	}
}

type CloudflareBypass struct {
	client *http.Client
}

func NewCloudflareBypass() *CloudflareBypass {
	return &CloudflareBypass{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CloudflareBypass) BypassChallenge(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 503 || resp.StatusCode == 403 {
		time.Sleep(5 * time.Second)
		return c.client.Do(req)
	}

	return resp, nil
}

type SessionManager struct {
	sessions map[string]*http.Client
	cookies  map[string][]*http.Cookie
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*http.Client),
		cookies:  make(map[string][]*http.Cookie),
	}
}

func (s *SessionManager) GetSession(domain string) *http.Client {
	if client, exists := s.sessions[domain]; exists {
		return client
	}

	jar := &cookieJar{cookies: make(map[string][]*http.Cookie)}
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	s.sessions[domain] = client
	return client
}

type cookieJar struct {
	cookies map[string][]*http.Cookie
}

func (j *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.cookies[u.Host] = cookies
}

func (j *cookieJar) Cookies(u *url.URL) []*http.Cookie {
	return j.cookies[u.Host]
}

type BotDetectionEvasion struct {
	stealthClient *StealthClient
	cfBypass      *CloudflareBypass
	sessionMgr    *SessionManager
}

func NewBotDetectionEvasion() *BotDetectionEvasion {
	config := &StealthConfig{
		RotateUserAgents:  true,
		RandomizeHeaders:  true,
		SimulateHuman:     true,
		BypassCloudflare:  true,
		DelayRange:        [2]int{1000, 5000},
		MaxRetries:        3,
		TLSFingerprinting: true,
	}

	return &BotDetectionEvasion{
		stealthClient: NewStealthClient(config),
		cfBypass:      NewCloudflareBypass(),
		sessionMgr:    NewSessionManager(),
	}
}

func (b *BotDetectionEvasion) MakeRequest(url string) (*http.Response, error) {
	domain := extractDomain(url)
	client := b.sessionMgr.GetSession(domain)

	req, err := b.stealthClient.CreateStealthRequest("GET", url)
	if err != nil {
		return nil, err
	}

	b.stealthClient.SimulateHumanDelay()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if isBlocked(resp) {
		return b.cfBypass.BypassChallenge(url)
	}

	return resp, nil
}

func isBlocked(resp *http.Response) bool {
	return resp.StatusCode == 403 || resp.StatusCode == 503 || 
		   resp.StatusCode == 429 || resp.StatusCode == 520
}

func extractDomain(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return url
}