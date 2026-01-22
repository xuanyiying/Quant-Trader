package matching

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAbuseRiskManager_RateLimit(t *testing.T) {
	engine := NewEngine()
	symbol := "RISK-USDT"
	engine.AddOrderBook(symbol)
	_ = engine.ConfigureSymbol(symbol, SymbolConfig{
		TickSize: decimal.RequireFromString("0.00000001"),
		LotSize:  decimal.RequireFromString("0.00000001"),
	})

	rm := NewAbuseRiskManager(AbuseConfig{
		MaxOrdersPerSecond: 1,
	})
	assert.NoError(t, engine.SetRiskManager(symbol, rm))

	o1 := NewOrder("r1", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	o1.AccountID = "A"
	_, err := engine.ProcessOrder(o1)
	assert.NoError(t, err)

	o2 := NewOrder("r2", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	o2.AccountID = "A"
	_, err = engine.ProcessOrder(o2)
	assert.Error(t, err)
}

func TestAbuseRiskManager_SpoofingBlock(t *testing.T) {
	engine := NewEngine()
	symbol := "SPOOF-USDT"
	engine.AddOrderBook(symbol)
	_ = engine.ConfigureSymbol(symbol, SymbolConfig{
		TickSize: decimal.RequireFromString("0.00000001"),
		LotSize:  decimal.RequireFromString("0.00000001"),
	})

	rm := NewAbuseRiskManager(AbuseConfig{
		MaxOrdersPerSecond: 10000,
		SpoofNotional:      decimal.RequireFromString("1"),
		SpoofWindowMs:      60_000,
		BlockDurationMs:    60_000,
	})
	assert.NoError(t, engine.SetRiskManager(symbol, rm))

	o1 := NewOrder("s1", symbol, SideSell, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	o1.AccountID = "A"
	_, err := engine.ProcessOrder(o1)
	assert.NoError(t, err)

	_, err = engine.CancelOrder(symbol, "s1")
	assert.NoError(t, err)

	o2 := NewOrder("s2", symbol, SideSell, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	o2.AccountID = "A"
	_, err = engine.ProcessOrder(o2)
	assert.Error(t, err)
}

func TestAbuseRiskManager_BlockSelfTrade(t *testing.T) {
	engine := NewEngine()
	symbol := "SELF-USDT"
	engine.AddOrderBook(symbol)
	_ = engine.ConfigureSymbol(symbol, SymbolConfig{
		TickSize: decimal.RequireFromString("0.00000001"),
		LotSize:  decimal.RequireFromString("0.00000001"),
	})

	rm := NewAbuseRiskManager(AbuseConfig{
		MaxOrdersPerSecond: 10000,
		BlockSelfTrade:     true,
		BlockDurationMs:    60_000,
	})
	assert.NoError(t, engine.SetRiskManager(symbol, rm))

	sell := NewOrder("ss1", symbol, SideSell, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	sell.AccountID = "A"
	_, err := engine.ProcessOrder(sell)
	assert.NoError(t, err)

	buy := NewOrder("sb1", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	buy.AccountID = "A"
	_, err = engine.ProcessOrder(buy)
	assert.NoError(t, err)

	o2 := NewOrder("sb2", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	o2.AccountID = "A"
	_, err = engine.ProcessOrder(o2)
	assert.Error(t, err)
}
