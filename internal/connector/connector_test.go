package connector

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBinanceConnector_ConvertToModel(t *testing.T) {
	logger := zap.NewNop()
	c := NewBinanceConnector(logger, "btcusdt")

	event := BinanceTradeEvent{
		TradeID:      12345,
		Price:        "50000.00",
		Quantity:     "0.1",
		TradeTime:    1640123456789,
		Symbol:       "BTCUSDT",
		IsBuyerMaker: true,
	}

	trade := c.convertToModel(event)

	assert.Equal(t, "12345", trade.ID)
	assert.Equal(t, "BTCUSDT", trade.Symbol)
	assert.Equal(t, "binance", trade.Exchange)
	assert.True(t, trade.Price.Equal(decimal.NewFromFloat(50000.00)))
	assert.True(t, trade.Amount.Equal(decimal.NewFromFloat(0.1)))
	assert.Equal(t, "sell", trade.Side) // IsBuyerMaker=true means sell
	assert.Equal(t, time.Unix(0, 1640123456789*int64(time.Millisecond)), trade.Timestamp)
}

func TestOKXConnector_ConvertToModel(t *testing.T) {
	logger := zap.NewNop()
	c := NewOKXConnector(logger, "BTC-USDT")

	data := OKXTradeData{
		TradeId: "98765",
		Px:      "50100.5",
		Sz:      "0.05",
		Side:    "buy",
		Ts:      "1640123456789",
		InstId:  "BTC-USDT",
	}

	trade := c.convertToModel(data)

	assert.Equal(t, "98765", trade.ID)
	assert.Equal(t, "BTC-USDT", trade.Symbol)
	assert.Equal(t, "okx", trade.Exchange)
	assert.True(t, trade.Price.Equal(decimal.NewFromFloat(50100.5)))
	assert.True(t, trade.Amount.Equal(decimal.NewFromFloat(0.05)))
	assert.Equal(t, "buy", trade.Side)
}

func TestBybitConnector_ConvertToModel(t *testing.T) {
	logger := zap.NewNop()
	c := NewBybitConnector(logger, "BTCUSDT")

	data := BybitTradeData{
		I:  "654321",
		S:  "BTCUSDT",
		S2: "Sell",
		P:  "49999.9",
		V:  "1.2",
		T:  1640123456789,
	}

	trade := c.convertToModel(data)

	assert.Equal(t, "654321", trade.ID)
	assert.Equal(t, "bybit", trade.Exchange)
	assert.Equal(t, "sell", trade.Side)
	assert.True(t, trade.Price.Equal(decimal.NewFromFloat(49999.9)))
}

func TestKrakenConnector_ConvertToModel(t *testing.T) {
	logger := zap.NewNop()
	c := NewKrakenConnector(logger, "XBT/USD")

	// Kraken data format: [price, volume, time, side, orderType, misc]
	data := []interface{}{
		"50000.1",
		"0.5",
		"1640123456.7890",
		"s",
		"m",
		"",
	}

	trade := c.convertToModel(data, "XBT/USD")

	assert.Equal(t, "kraken", trade.Exchange)
	assert.Equal(t, "sell", trade.Side)
	assert.True(t, trade.Price.Equal(decimal.NewFromFloat(50000.1)))
	assert.True(t, trade.Amount.Equal(decimal.NewFromFloat(0.5)))
}

func TestCoinbaseConnector_ConvertToModel(t *testing.T) {
	logger := zap.NewNop()
	c := NewCoinbaseConnector(logger, "BTC-USD")

	event := CoinbaseMatchEvent{
		TradeID:   555,
		ProductID: "BTC-USD",
		Price:     "50050.00",
		Size:      "0.25",
		Side:      "buy",
		Time:      "2021-12-22T00:00:00Z",
	}

	trade := c.convertToModel(event)

	assert.Equal(t, "555", trade.ID)
	assert.Equal(t, "coinbase", trade.Exchange)
	assert.True(t, trade.Price.Equal(decimal.NewFromFloat(50050.00)))
	assert.Equal(t, "buy", trade.Side)
}
