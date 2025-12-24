package engine

import (
	"context"
	"quant-trader/internal/model"
	"quant-trader/internal/strategy"

	"go.uber.org/zap"
)

type WorkerPool struct {
	jobQueue    chan model.Trade
	workerCount int
	algo        strategy.Strategy
	logger      *zap.Logger
}

func NewWorkerPool(workerCount int, bufferSize int, algo strategy.Strategy, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan model.Trade, bufferSize),
		workerCount: workerCount,
		algo:        algo,
		logger:      logger,
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
	// In a real scenario, we might want to aggregate data before calling the strategy
	// Or call OnData if the strategy supports it (we added it to MAStrategy)

	// For now, just a placeholder for real-time strategy execution
	p.logger.Debug("worker processing trade",
		zap.Int("worker_id", workerID),
		zap.String("symbol", trade.Symbol),
		zap.String("price", trade.Price.String()),
	)
}
