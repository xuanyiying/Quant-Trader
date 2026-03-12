package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Metrics 封装所有 Prometheus 指标
type Metrics struct {
	logger *zap.Logger

	// HTTP 指标
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec

	// 业务指标
	ActiveUsers      prometheus.Gauge
	PaperOrdersTotal *prometheus.CounterVec
	PaperTradesTotal *prometheus.CounterVec

	// K线数据指标
	KlinesInserted   prometheus.Counter
	KlinesBackfilled *prometheus.CounterVec

	// 回测指标
	BacktestsRunning prometheus.Gauge
	BacktestsTotal   *prometheus.CounterVec

	// 策略市场指标
	StrategiesPublished prometheus.Gauge
	StrategiesPurchased *prometheus.CounterVec

	// 系统指标
	Goroutines prometheus.Gauge
	MemoryUsage prometheus.Gauge
}

// NewMetrics 创建指标收集器
func NewMetrics(logger *zap.Logger) *Metrics {
	m := &Metrics{
		logger: logger,
	}

	// HTTP 请求计数
	m.HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP 请求耗时
	m.HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// 活跃用户
	m.ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of active users",
		},
	)

	// 模拟订单
	m.PaperOrdersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paper_orders_total",
			Help: "Total number of paper orders",
		},
		[]string{"symbol", "side", "type"},
	)

	// 模拟交易
	m.PaperTradesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paper_trades_total",
			Help: "Total number of paper trades",
		},
		[]string{"symbol", "side"},
	)

	// K线插入
	m.KlinesInserted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "klines_inserted_total",
			Help: "Total number of klines inserted",
		},
	)

	// K线补全
	m.KlinesBackfilled = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "klines_backfilled_total",
			Help: "Total number of klines backfilled",
		},
		[]string{"exchange", "symbol"},
	)

	// 回测运行中
	m.BacktestsRunning = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "backtests_running",
			Help: "Number of backtests currently running",
		},
	)

	// 回测总数
	m.BacktestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "backtests_total",
			Help: "Total number of backtests",
		},
		[]string{"strategy_type", "status"},
	)

	// 策略发布
	m.StrategiesPublished = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "strategies_published",
			Help: "Number of strategies published in marketplace",
		},
	)

	// 策略购买
	m.StrategiesPurchased = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "strategies_purchased_total",
			Help: "Total number of strategies purchased",
		},
		[]string{"strategy_id"},
	)

	// Goroutines
	m.Goroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_goroutines",
			Help: "Number of goroutines",
		},
	)

	// 内存使用
	m.MemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	return m
}

// RecordHTTPRequest 记录 HTTP 请求
func (m *Metrics) RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordPaperOrder 记录模拟订单
func (m *Metrics) RecordPaperOrder(symbol, side, orderType string) {
	m.PaperOrdersTotal.WithLabelValues(symbol, side, orderType).Inc()
}

// RecordPaperTrade 记录模拟交易
func (m *Metrics) RecordPaperTrade(symbol, side string) {
	m.PaperTradesTotal.WithLabelValues(symbol, side).Inc()
}

// RecordKlinesInserted 记录 K线插入
func (m *Metrics) RecordKlinesInserted(count int) {
	m.KlinesInserted.Add(float64(count))
}

// RecordKlinesBackfilled 记录 K线补全
func (m *Metrics) RecordKlinesBackfilled(exchange, symbol string, count int) {
	m.KlinesBackfilled.WithLabelValues(exchange, symbol).Add(float64(count))
}

// RecordBacktestStart 记录回测开始
func (m *Metrics) RecordBacktestStart(strategyType string) {
	m.BacktestsRunning.Inc()
}

// RecordBacktestEnd 记录回测结束
func (m *Metrics) RecordBacktestEnd(strategyType, status string) {
	m.BacktestsRunning.Dec()
	m.BacktestsTotal.WithLabelValues(strategyType, status).Inc()
}

// RecordStrategyPurchase 记录策略购买
func (m *Metrics) RecordStrategyPurchase(strategyID string) {
	m.StrategiesPurchased.WithLabelValues(strategyID).Inc()
}

// UpdateActiveUsers 更新活跃用户数
func (m *Metrics) UpdateActiveUsers(count int) {
	m.ActiveUsers.Set(float64(count))
}

// UpdateStrategiesPublished 更新已发布策略数
func (m *Metrics) UpdateStrategiesPublished(count int) {
	m.StrategiesPublished.Set(float64(count))
}

// StartMetricsServer 启动指标服务器
func (m *Metrics) StartMetricsServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	m.logger.Info("starting metrics server", zap.String("addr", addr))
	return http.ListenAndServe(addr, mux)
}

// Middleware HTTP 中间件
func (m *Metrics) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 包装 ResponseWriter 以获取状态码
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			m.RecordHTTPRequest(r.Method, r.URL.Path, http.StatusText(wrapped.statusCode), duration)
		})
	}
}

// responseWriter 包装 http.ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
