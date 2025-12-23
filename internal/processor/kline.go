package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"market-ingestor/internal/infrastructure"
	"market-ingestor/internal/model"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type KlineProcessor struct {
	js      nats.JetStreamContext
	logger  *zap.Logger
	candles map[string]*model.KLine
	mu      sync.Mutex
}

func NewKlineProcessor(js nats.JetStreamContext, logger *zap.Logger) *KlineProcessor {
	return &KlineProcessor{
		js:      js,
		logger:  logger,
		candles: make(map[string]*model.KLine),
	}
}

func (p *KlineProcessor) Run(ctx context.Context) error {
	_, err := p.js.Subscribe("market.raw.*.*", func(msg *nats.Msg) {
		var trade model.Trade
		if err := json.Unmarshal(msg.Data, &trade); err != nil {
			p.logger.Error("failed to unmarshal trade in processor", zap.Error(err))
			return
		}
		infrastructure.TradeProcessRate.WithLabelValues(trade.Symbol).Inc()
		p.processTrade(trade)
		msg.Ack()
	}, nats.Durable("kline-processor"), nats.ManualAck())

	if err != nil {
		return err
	}

	go p.flushLoop(ctx)
	p.logger.Info("kline processor started")
	return nil
}

func (p *KlineProcessor) processTrade(trade model.Trade) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Use 1 minute resolution
	window := trade.Timestamp.Truncate(time.Minute)
	key := fmt.Sprintf("%s:%s:%s", trade.Exchange, trade.Symbol, window.Format(time.RFC3339))

	candle, ok := p.candles[key]
	if !ok {
		candle = &model.KLine{
			Symbol:    trade.Symbol,
			Exchange:  trade.Exchange,
			Period:    "1m",
			Open:      trade.Price,
			High:      trade.Price,
			Low:       trade.Price,
			Close:     trade.Price,
			Volume:    trade.Amount,
			Timestamp: window,
		}
		p.candles[key] = candle
	} else {
		if trade.Price.GreaterThan(candle.High) {
			candle.High = trade.Price
		}
		if trade.Price.LessThan(candle.Low) {
			candle.Low = trade.Price
		}
		candle.Close = trade.Price
		candle.Volume = candle.Volume.Add(trade.Amount)
	}
}

func (p *KlineProcessor) flushLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.flush()
		}
	}
}

func (p *KlineProcessor) flush() {
	p.mu.Lock()
	now := time.Now().Truncate(time.Minute)
	toFlush := make([]*model.KLine, 0)

	for key, candle := range p.candles {
		// If candle timestamp is before current minute, it's completed
		if candle.Timestamp.Before(now) {
			toFlush = append(toFlush, candle)
			delete(p.candles, key)
		}
	}
	p.mu.Unlock()

	for _, candle := range toFlush {
		subject := fmt.Sprintf("market.kline.1m.%s", candle.Symbol)
		data, _ := json.Marshal(candle)
		_, err := p.js.Publish(subject, data)
		if err != nil {
			p.logger.Error("failed to publish kline", zap.Error(err))
		}
	}
}
