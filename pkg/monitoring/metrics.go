package monitoring

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Metrics struct {
	RequestsTotal     *prometheus.CounterVec
	RequestDuration   *prometheus.HistogramVec
	RequestsInFlight  *prometheus.GaugeVec
	
	ResponseSize      *prometheus.HistogramVec
	ResponseStatus    *prometheus.CounterVec
	
	CacheHits         *prometheus.CounterVec
	CacheMisses       *prometheus.CounterVec
	CacheSize         *prometheus.GaugeVec
	
	QueueSize         *prometheus.GaugeVec
	QueueProcessed    *prometheus.CounterVec
	QueueErrors       *prometheus.CounterVec
	
	BrowserSessions   *prometheus.GaugeVec
	BrowserErrors     *prometheus.CounterVec
	PageLoadTime      *prometheus.HistogramVec
	
	MemoryUsage       *prometheus.GaugeVec
	CPUUsage          *prometheus.GaugeVec
	GoroutineCount    prometheus.Gauge
	
	DataExtracted     *prometheus.CounterVec
	ErrorsTotal       *prometheus.CounterVec
	RetryAttempts     *prometheus.CounterVec
	
	registry *prometheus.Registry
	logger   *zap.Logger
}

func NewMetrics(logger *zap.Logger) *Metrics {
	registry := prometheus.NewRegistry()
	
	m := &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "status", "host"},
		),
		
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "goscraper_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "host"},
		),
		
		RequestsInFlight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goscraper_requests_in_flight",
				Help: "Number of HTTP requests currently in flight",
			},
			[]string{"host"},
		),
		
		ResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "goscraper_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: []float64{1024, 4096, 16384, 65536, 262144, 1048576, 4194304},
			},
			[]string{"host"},
		),
		
		ResponseStatus: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_response_status_total",
				Help: "Total number of responses by status code",
			},
			[]string{"status", "host"},
		),
		
		CacheHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_type"},
		),
		
		CacheMisses: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_type"},
		),
		
		CacheSize: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goscraper_cache_size_bytes",
				Help: "Current cache size in bytes",
			},
			[]string{"cache_type"},
		),
		
		QueueSize: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goscraper_queue_size",
				Help: "Current queue size",
			},
			[]string{"queue_name", "priority"},
		),
		
		QueueProcessed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_queue_processed_total",
				Help: "Total number of processed queue items",
			},
			[]string{"queue_name", "status"},
		),
		
		QueueErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_queue_errors_total",
				Help: "Total number of queue processing errors",
			},
			[]string{"queue_name", "error_type"},
		),
		
		BrowserSessions: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goscraper_browser_sessions",
				Help: "Number of active browser sessions",
			},
			[]string{"engine"},
		),
		
		BrowserErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_browser_errors_total",
				Help: "Total number of browser errors",
			},
			[]string{"engine", "error_type"},
		),
		
		PageLoadTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "goscraper_page_load_time_seconds",
				Help:    "Page load time in seconds",
				Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
			},
			[]string{"engine", "host"},
		),
		
		MemoryUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goscraper_memory_usage_bytes",
				Help: "Memory usage in bytes",
			},
			[]string{"type"},
		),
		
		CPUUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goscraper_cpu_usage_percent",
				Help: "CPU usage percentage",
			},
			[]string{"core"},
		),
		
		GoroutineCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "goscraper_goroutines",
				Help: "Number of goroutines",
			},
		),
		
		DataExtracted: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_data_extracted_total",
				Help: "Total amount of data extracted",
			},
			[]string{"type", "source"},
		),
		
		ErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_errors_total",
				Help: "Total number of errors",
			},
			[]string{"type", "component"},
		),
		
		RetryAttempts: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goscraper_retry_attempts_total",
				Help: "Total number of retry attempts",
			},
			[]string{"component", "reason"},
		),
		
		registry: registry,
		logger:   logger,
	}
	
	m.registerMetrics()
	
	return m
}

func (m *Metrics) registerMetrics() {
	m.registry.MustRegister(
		m.RequestsTotal,
		m.RequestDuration,
		m.RequestsInFlight,
		m.ResponseSize,
		m.ResponseStatus,
		m.CacheHits,
		m.CacheMisses,
		m.CacheSize,
		m.QueueSize,
		m.QueueProcessed,
		m.QueueErrors,
		m.BrowserSessions,
		m.BrowserErrors,
		m.PageLoadTime,
		m.MemoryUsage,
		m.CPUUsage,
		m.GoroutineCount,
		m.DataExtracted,
		m.ErrorsTotal,
		m.RetryAttempts,
	)
}

func (m *Metrics) RecordRequest(method, host, status string, duration time.Duration, size int64) {
	m.RequestsTotal.WithLabelValues(method, status, host).Inc()
	m.RequestDuration.WithLabelValues(method, host).Observe(duration.Seconds())
	m.ResponseSize.WithLabelValues(host).Observe(float64(size))
	m.ResponseStatus.WithLabelValues(status, host).Inc()
}

func (m *Metrics) RecordCacheHit(cacheType string) {
	m.CacheHits.WithLabelValues(cacheType).Inc()
}

func (m *Metrics) RecordCacheMiss(cacheType string) {
	m.CacheMisses.WithLabelValues(cacheType).Inc()
}

func (m *Metrics) RecordQueueSize(queueName, priority string, size float64) {
	m.QueueSize.WithLabelValues(queueName, priority).Set(size)
}

func (m *Metrics) RecordBrowserSession(engine string, delta float64) {
	m.BrowserSessions.WithLabelValues(engine).Add(delta)
}

func (m *Metrics) RecordPageLoad(engine, host string, duration time.Duration) {
	m.PageLoadTime.WithLabelValues(engine, host).Observe(duration.Seconds())
}

func (m *Metrics) RecordError(errorType, component string) {
	m.ErrorsTotal.WithLabelValues(errorType, component).Inc()
}

func (m *Metrics) RecordRetry(component, reason string) {
	m.RetryAttempts.WithLabelValues(component, reason).Inc()
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *Metrics) StartMetricsServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", m.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	
	m.logger.Info("Starting metrics server", zap.String("addr", addr))
	return server.ListenAndServe()
}

type Alert struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Query       string            `json:"query"`
	Threshold   float64           `json:"threshold"`
	Duration    time.Duration     `json:"duration"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type AlertManager struct {
	alerts  map[string]*Alert
	metrics *Metrics
	logger  *zap.Logger
}

func NewAlertManager(metrics *Metrics, logger *zap.Logger) *AlertManager {
	return &AlertManager{
		alerts:  make(map[string]*Alert),
		metrics: metrics,
		logger:  logger,
	}
}

func (a *AlertManager) AddAlert(alert *Alert) {
	a.alerts[alert.Name] = alert
}

func (a *AlertManager) CheckAlerts(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for name, alert := range a.alerts {
				if a.evaluateAlert(alert) {
					a.logger.Warn("Alert triggered",
						zap.String("alert", name),
						zap.String("description", alert.Description),
					)
					//TODO: NOTIFICATION SYSTEM ENTEGRATION
				}
			}
		}
	}
}

func (a *AlertManager) evaluateAlert(alert *Alert) bool {
	
	return false
}