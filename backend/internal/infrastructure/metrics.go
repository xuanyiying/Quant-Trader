package infrastructure

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	IngestLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "ingest_latency_seconds",
		Help: "Latency of market data ingestion",
	}, []string{"exchange", "symbol"})

	WSConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ws_connections_total",
		Help: "Total number of active WebSocket connections",
	})

	DBInsertRate = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_insert_total",
		Help: "Total number of records inserted into DB",
	}, []string{"table"})

	TradeProcessRate = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "trade_process_total",
		Help: "Total number of trades processed",
	}, []string{"symbol"})

	GoroutineCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "goroutine_count",
		Help: "Number of active goroutines",
	})
)
