package matching

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func BenchmarkMatching_LimitTaker(b *testing.B) {
	engine := NewEngine()
	symbol := "BENCH-USDT"
	engine.AddOrderBook(symbol)
	_ = engine.ConfigureSymbol(symbol, SymbolConfig{
		TickSize: decimal.RequireFromString("0.00000001"),
		LotSize:  decimal.RequireFromString("0.00000001"),
	})
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-engine.Events():
			case <-stop:
				return
			}
		}
	}()
	defer close(stop)

	for i := 0; i < 2000; i++ {
		price := decimal.NewFromInt(1000).Add(decimal.NewFromInt(int64(i)))
		o := NewOrder(fmt.Sprintf("m-%d", i), symbol, SideSell, TypeLimit, price, decimal.RequireFromString("1"))
		_, _ = engine.ProcessOrder(o)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := NewOrder(fmt.Sprintf("t-%d", i), symbol, SideBuy, TypeLimit, decimal.NewFromInt(3000), decimal.RequireFromString("0.1"))
		_, _ = engine.ProcessOrder(o)
	}
}

func BenchmarkMatching_Cancel(b *testing.B) {
	engine := NewEngine()
	symbol := "BENCHC-USDT"
	engine.AddOrderBook(symbol)
	_ = engine.ConfigureSymbol(symbol, SymbolConfig{
		TickSize: decimal.RequireFromString("0.00000001"),
		LotSize:  decimal.RequireFromString("0.00000001"),
	})
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-engine.Events():
			case <-stop:
				return
			}
		}
	}()
	defer close(stop)

	orderIDs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("c-%d", i)
		orderIDs[i] = id
		o := NewOrder(id, symbol, SideBuy, TypeLimit, decimal.NewFromInt(1000), decimal.RequireFromString("0.1"))
		_, _ = engine.ProcessOrder(o)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.CancelOrder(symbol, orderIDs[i])
	}
}
