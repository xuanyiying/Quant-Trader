package strategy

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMAStrategy_Logic(t *testing.T) {
	s := NewMAStrategy(2, 4)

	// Not enough data
	prices := []int64{10, 11, 12}
	for _, p := range prices {
		s.OnData(decimal.NewFromInt(p))
	}
	// Action is only via OnCandle, but we can check calculateMA
	assert.Equal(t, 3, len(s.prices))

	s.OnData(decimal.NewFromInt(13))
	assert.Equal(t, 4, len(s.prices))

	// Short MA (2): (12+13)/2 = 12.5
	// Long MA (4): (10+11+12+13)/4 = 11.5
	// Short > Long -> Buy
	shortMA := s.calculateMA(2)
	longMA := s.calculateMA(4)
	assert.True(t, shortMA.GreaterThan(longMA))

	// Next price drops significantly
	s.OnData(decimal.NewFromInt(5))
	// Short MA (2): (13+5)/2 = 9
	// Long MA (4): (11+12+13+5)/4 = 10.25
	// Short < Long -> Sell
	shortMA = s.calculateMA(2)
	longMA = s.calculateMA(4)
	assert.True(t, shortMA.LessThan(longMA))
}
