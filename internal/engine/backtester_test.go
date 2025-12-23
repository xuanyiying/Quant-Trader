package engine

import (
	"market-ingestor/internal/model"
	"market-ingestor/internal/strategy"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestBacktester(t *testing.T) {
	// 1. Setup Strategy
	strat := strategy.NewMAStrategy(2, 5)
	initialBalance := decimal.NewFromInt(10000)
	tester := NewBacktester(strat, initialBalance)

	// 2. Generate dummy candles
	// A simple uptrend followed by a downtrend
	prices := []float64{100, 102, 104, 106, 108, 110, 108, 106, 104, 102, 100}
	candles := make([]model.KLine, len(prices))
	now := time.Now()

	for i, p := range prices {
		candles[i] = model.KLine{
			Symbol:    "BTCUSDT",
			Close:     decimal.NewFromFloat(p),
			Timestamp: now.Add(time.Duration(i) * time.Minute),
		}
	}

	// 3. Run Backtest
	report := tester.Run(candles)

	// 4. Verify
	if report.TotalTrades == 0 {
		t.Log("No trades were executed. Check strategy logic.")
	}

	if report.FinalBalance.Equal(initialBalance) && report.TotalTrades > 0 {
		t.Errorf("Final balance should be different if trades were made")
	}

	t.Logf("Backtest Report: %+v", report)
	for _, trade := range report.TradesLog {
		t.Logf("Trade: %s %s @ %s, Size: %s", trade.Side, trade.Symbol, trade.Price, trade.Size)
	}
}
