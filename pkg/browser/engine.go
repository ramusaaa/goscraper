package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Engine interface {
	Navigate(ctx context.Context, url string) error
	ExecuteScript(ctx context.Context, script string) (interface{}, error)
	Screenshot(ctx context.Context) ([]byte, error)
	GetHTML(ctx context.Context) (string, error)
	WaitForSelector(ctx context.Context, selector string, timeout time.Duration) error
	Click(ctx context.Context, selector string) error
	Type(ctx context.Context, selector, text string) error
	Close() error
}

type EngineType string

const (
	ChromeDP EngineType = "chromedp"
	Rod      EngineType = "rod"
)

type Config struct {
	Engine          EngineType
	Headless        bool
	UserAgent       string
	ViewportWidth   int
	ViewportHeight  int
	Timeout         time.Duration
	ProxyURL        string
	DisableImages   bool
	DisableCSS      bool
	DisableJS       bool
	CustomFlags     []string
	Extensions      []string
}

type Manager struct {
	config  *Config
	engines map[string]Engine
	pool    chan Engine
}

func NewManager(config *Config, poolSize int) *Manager {
	return &Manager{
		config:  config,
		engines: make(map[string]Engine),
		pool:    make(chan Engine, poolSize),
	}
}

func (m *Manager) GetEngine(ctx context.Context) (Engine, error) {
	select {
	case engine := <-m.pool:
		return engine, nil
	default:
		return m.createEngine(ctx)
	}
}

func (m *Manager) ReturnEngine(engine Engine) {
	select {
	case m.pool <- engine:
	default:
		engine.Close()
	}
}

func (m *Manager) createEngine(ctx context.Context) (Engine, error) {
	switch m.config.Engine {
	case ChromeDP:
		return m.createChromeDPEngine(ctx)
	case Rod:
		return m.createRodEngine(ctx)
	default:
		return nil, fmt.Errorf("unsupported engine: %s", m.config.Engine)
	}
}

type ChromeDPEngine struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (m *Manager) createChromeDPEngine(ctx context.Context) (*ChromeDPEngine, error) {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", m.config.Headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent(m.config.UserAgent),
		chromedp.WindowSize(m.config.ViewportWidth, m.config.ViewportHeight),
	}

	if m.config.ProxyURL != "" {
		opts = append(opts, chromedp.ProxyServer(m.config.ProxyURL))
	}

	if m.config.DisableImages {
		opts = append(opts, chromedp.Flag("blink-settings", "imagesEnabled=false"))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	engineCtx, _ := chromedp.NewContext(allocCtx)

	return &ChromeDPEngine{
		ctx:    engineCtx,
		cancel: cancel,
	}, nil
}

func (e *ChromeDPEngine) Navigate(ctx context.Context, url string) error {
	return chromedp.Run(e.ctx, chromedp.Navigate(url))
}

func (e *ChromeDPEngine) ExecuteScript(ctx context.Context, script string) (interface{}, error) {
	var result interface{}
	err := chromedp.Run(e.ctx, chromedp.Evaluate(script, &result))
	return result, err
}

func (e *ChromeDPEngine) Screenshot(ctx context.Context) ([]byte, error) {
	var buf []byte
	err := chromedp.Run(e.ctx, chromedp.CaptureScreenshot(&buf))
	return buf, err
}

func (e *ChromeDPEngine) GetHTML(ctx context.Context) (string, error) {
	var html string
	err := chromedp.Run(e.ctx, chromedp.OuterHTML("html", &html))
	return html, err
}

func (e *ChromeDPEngine) WaitForSelector(ctx context.Context, selector string, timeout time.Duration) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return chromedp.Run(timeoutCtx, chromedp.WaitVisible(selector))
}

func (e *ChromeDPEngine) Click(ctx context.Context, selector string) error {
	return chromedp.Run(e.ctx, chromedp.Click(selector))
}

func (e *ChromeDPEngine) Type(ctx context.Context, selector, text string) error {
	return chromedp.Run(e.ctx, chromedp.SendKeys(selector, text))
}

func (e *ChromeDPEngine) Close() error {
	e.cancel()
	return nil
}

type RodEngine struct {
	browser *rod.Browser
	page    *rod.Page
}

func (m *Manager) createRodEngine(ctx context.Context) (*RodEngine, error) {
	browser := rod.New()
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	page := browser.MustPage()

	return &RodEngine{
		browser: browser,
		page:    page,
	}, nil
}

func (e *RodEngine) Navigate(ctx context.Context, url string) error {
	return e.page.Navigate(url)
}

func (e *RodEngine) ExecuteScript(ctx context.Context, script string) (interface{}, error) {
	result, err := e.page.Eval(script)
	if err != nil {
		return nil, err
	}
	return result.Value, nil
}

func (e *RodEngine) Screenshot(ctx context.Context) ([]byte, error) {
	return e.page.Screenshot(true, nil)
}

func (e *RodEngine) GetHTML(ctx context.Context) (string, error) {
	return e.page.HTML()
}

func (e *RodEngine) WaitForSelector(ctx context.Context, selector string, timeout time.Duration) error {
	element, err := e.page.Timeout(timeout).Element(selector)
	if err != nil {
		return err
	}
	return element.WaitVisible()
}

func (e *RodEngine) Click(ctx context.Context, selector string) error {
	element, err := e.page.Element(selector)
	if err != nil {
		return err
	}
	return element.Click(proto.InputMouseButtonLeft, 1)
}

func (e *RodEngine) Type(ctx context.Context, selector, text string) error {
	element, err := e.page.Element(selector)
	if err != nil {
		return err
	}
	return element.Input(text)
}

func (e *RodEngine) Close() error {
	if e.page != nil {
		e.page.Close()
	}
	if e.browser != nil {
		e.browser.Close()
	}
	return nil
}