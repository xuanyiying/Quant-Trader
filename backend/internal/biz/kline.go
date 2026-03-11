package biz

import (
	"context"
	"time"

	"quant-trader/internal/engine"
	"quant-trader/internal/model"
	"quant-trader/internal/repo"
	"quant-trader/internal/strategy"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Kline struct {
	repo   *repo.Kline
	logger *zap.Logger
}

func NewKline(repo *repo.Kline, logger *zap.Logger) *Kline {
	return &Kline{repo: repo, logger: logger}
}

func (b *Kline) GetLatest(ctx context.Context, symbol, period string, limit int) ([]model.KLine, error) {
	return b.repo.GetLatest(ctx, symbol, period, limit)
}

func (b *Kline) GetBySymbol(ctx context.Context, symbol, period string, start, end time.Time) ([]model.KLine, error) {
	return b.repo.GetBySymbol(ctx, symbol, period, start, end, 0)
}

type BacktestInput struct {
	Symbol         string
	StrategyType   string
	Config         map[string]interface{}
	InitialBalance decimal.Decimal
	StartTime      time.Time
	EndTime        time.Time
}

func (b *Kline) RunBacktest(ctx context.Context, input BacktestInput) (interface{}, error) {
	klines, err := b.repo.GetBySymbol(ctx, input.Symbol, model.Period1m, input.StartTime, input.EndTime, 0)
	if err != nil {
		return nil, err
	}

	strat, err := strategy.NewStrategy(input.StrategyType, input.Config)
	if err != nil {
		return nil, err
	}

	tester := engine.NewBacktester(strat, input.InitialBalance)
	return tester.Run(klines), nil
}
