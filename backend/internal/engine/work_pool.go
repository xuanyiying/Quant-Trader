package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"quant-trader/internal/model"
	"quant-trader/internal/strategy"
	"strings"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type WorkerPool struct {
	jobQueue    chan model.Trade
	workerCount int
	algo        strategy.Strategy
	logger      *zap.Logger
	js          nats.JetStreamContext
}

func NewWorkerPool(workerCount int, bufferSize int, algo strategy.Strategy, logger *zap.Logger, js nats.JetStreamContext) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan model.Trade, bufferSize),
		workerCount: workerCount,
		algo:        algo,
		logger:      logger,
		js:          js,
	}
}

func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		go p.worker(ctx, i)
	}
	p.logger.Info("started worker pool", zap.Int("workers", p.workerCount))
}

func (p *WorkerPool) Submit(trade model.Trade) {
	select {
	case p.jobQueue <- trade:
	default:
		p.logger.Warn("worker pool job queue full, dropping trade")
	}
}

func (p *WorkerPool) worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		case trade, ok := <-p.jobQueue:
			if !ok {
				return
			}
			p.process(id, trade)
		}
	}
}

func (p *WorkerPool) process(workerID int, trade model.Trade) {
	action := p.algo.OnTrade(trade)

	if action != strategy.ActionHold {
		p.logger.Info("WORKER STRATEGY SIGNAL",
			zap.Int("worker_id", workerID),
			zap.String("strategy", p.algo.Name()),
			zap.String("symbol", trade.Symbol),
			zap.String("action", string(action)),
			zap.String("price", trade.Price.String()),
			zap.Time("time", trade.Timestamp),
		)

		// Publish signal to NATS
		signalSubject := fmt.Sprintf("%s.%s.%s", model.SubjectStrategySignal, p.algo.Name(), trade.Symbol)
		signalData := map[string]interface{}{
			"strategy": p.algo.Name(),
			"symbol":   trade.Symbol,
			"action":   strings.ToLower(string(action)),
			"price":    trade.Price.String(),
			"time":     trade.Timestamp,
			"source":   "worker_pool",
		}

		data, err := json.Marshal(signalData)
		if err != nil {
			p.logger.Error("failed to marshal signal data", zap.Error(err))
			return
		}

		if p.js != nil {
			_, err = p.js.Publish(signalSubject, data)
			if err != nil {
				p.logger.Error("failed to publish signal to NATS", zap.Error(err))
			}
		}
	}
}
