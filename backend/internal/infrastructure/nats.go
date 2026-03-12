package infrastructure

import (
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func InitNATS(url string, logger *zap.Logger) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, nil, err
	}

	// Create stream if it doesn't exist
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "MARKET",
		Subjects: []string{"market.raw.*.*", "market.kline.*.*"},
	})
	if err != nil {
		_, _ = js.UpdateStream(&nats.StreamConfig{
			Name:     "MARKET",
			Subjects: []string{"market.raw.*.*", "market.kline.*.*"},
		})
	}

	// Create STRATEGY stream for strategy signals
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "STRATEGY",
		Subjects:  []string{"strategy.signal.*.*"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxAge:    24 * time.Hour,
	})
	if err != nil {
		_, _ = js.UpdateStream(&nats.StreamConfig{
			Name:     "STRATEGY",
			Subjects: []string{"strategy.signal.*.*"},
		})
	}

	// Create PAPER stream for paper trading events
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "PAPER",
		Subjects:  []string{"paper.event.*"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxAge:    24 * time.Hour,
	})
	if err != nil {
		_, _ = js.UpdateStream(&nats.StreamConfig{
			Name:     "PAPER",
			Subjects: []string{"paper.event.*"},
		})
	}
	if err != nil {
		// If stream exists, we might need to update it
		_, err = js.UpdateStream(&nats.StreamConfig{
			Name:     "MARKET",
			Subjects: []string{"market.raw.*.*", "market.kline.*.*"},
		})
		if err != nil {
			logger.Warn("failed to create or update stream", zap.Error(err))
		}
	}

	return nc, js, nil
}
