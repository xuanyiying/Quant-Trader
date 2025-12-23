package processor

import (
	"market-ingestor/internal/model"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestKlineProcessor_ProcessTrade(t *testing.T) {
	logger := zap.NewNop()
	p := NewKlineProcessor(nil, logger)

	now := time.Now().Truncate(time.Minute)
	symbol := "BTCUSDT"
	exchange := "binance"

	// 1. First trade creates the candle
	trade1 := model.Trade{
		ID:        "1",
		Symbol:    symbol,
		Exchange:  exchange,
		Price:     decimal.NewFromFloat(50000),
		Amount:    decimal.NewFromFloat(1),
		Timestamp: now.Add(10 * time.Second),
	}
	p.processTrade(trade1)

	key := "binance:BTCUSDT:" + now.Format(time.RFC3339)
	candle, ok := p.candles[key]
	assert.True(t, ok)
	assert.True(t, candle.Open.Equal(decimal.NewFromFloat(50000)))
	assert.True(t, candle.High.Equal(decimal.NewFromFloat(50000)))
	assert.True(t, candle.Low.Equal(decimal.NewFromFloat(50000)))
	assert.True(t, candle.Close.Equal(decimal.NewFromFloat(50000)))
	assert.True(t, candle.Volume.Equal(decimal.NewFromFloat(1)))

	// 2. Second trade updates high and close
	trade2 := model.Trade{
		ID:        "2",
		Symbol:    symbol,
		Exchange:  exchange,
		Price:     decimal.NewFromFloat(50100),
		Amount:    decimal.NewFromFloat(0.5),
		Timestamp: now.Add(20 * time.Second),
	}
	p.processTrade(trade2)

	assert.True(t, candle.High.Equal(decimal.NewFromFloat(50100)))
	assert.True(t, candle.Low.Equal(decimal.NewFromFloat(50000)))
	assert.True(t, candle.Close.Equal(decimal.NewFromFloat(50100)))
	assert.True(t, candle.Volume.Equal(decimal.NewFromFloat(1.5)))

	// 3. Third trade updates low and close
	trade3 := model.Trade{
		ID:        "3",
		Symbol:    symbol,
		Exchange:  exchange,
		Price:     decimal.NewFromFloat(49900),
		Amount:    decimal.NewFromFloat(2),
		Timestamp: now.Add(30 * time.Second),
	}
	p.processTrade(trade3)

	assert.True(t, candle.High.Equal(decimal.NewFromFloat(50100)))
	assert.True(t, candle.Low.Equal(decimal.NewFromFloat(49900)))
	assert.True(t, candle.Close.Equal(decimal.NewFromFloat(49900)))
	assert.True(t, candle.Volume.Equal(decimal.NewFromFloat(3.5)))
}
