package engine

import (
	"context"
	"market-ingestor/internal/model"
	"market-ingestor/internal/strategy"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func TestWorkerPool_SubmitAndProcess(t *testing.T) {
	logger := zap.NewNop()
	strat := strategy.NewMAStrategy(5, 10)
	pool := NewWorkerPool(2, 10, strat, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	trade := model.Trade{
		Symbol: "BTCUSDT",
		Price:  decimal.NewFromInt(50000),
		Amount: decimal.NewFromFloat(1.5),
	}

	// Submit multiple trades
	for i := 0; i < 5; i++ {
		pool.Submit(trade)
	}

	// Allow some time for workers to pick up jobs
	time.Sleep(100 * time.Millisecond)
}
