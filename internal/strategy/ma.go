package strategy

import (
	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
)

type MAStrategy struct {
	shortPeriod int
	longPeriod  int
	prices      []decimal.Decimal
}

func NewMAStrategy(short, long int) *MAStrategy {
	return &MAStrategy{
		shortPeriod: short,
		longPeriod:  long,
		prices:      make([]decimal.Decimal, 0),
	}
}

func (s *MAStrategy) Name() string {
	return "Moving Average Crossover"
}

func (s *MAStrategy) OnData(price decimal.Decimal) {
	s.prices = append(s.prices, price)
	if len(s.prices) > s.longPeriod+1 {
		s.prices = s.prices[1:]
	}
}

func (s *MAStrategy) OnCandle(candle model.KLine) Action {
	s.OnData(candle.Close)

	if len(s.prices) < s.longPeriod {
		return ActionHold
	}

	shortMA := s.calculateMA(s.shortPeriod)
	longMA := s.calculateMA(s.longPeriod)

	// Simple crossover logic
	if shortMA.GreaterThan(longMA) {
		return ActionBuy
	} else if shortMA.LessThan(longMA) {
		return ActionSell
	}

	return ActionHold
}

func (s *MAStrategy) calculateMA(period int) decimal.Decimal {
	sum := decimal.Zero
	data := s.prices[len(s.prices)-period:]
	for _, p := range data {
		sum = sum.Add(p)
	}
	return sum.Div(decimal.NewFromInt(int64(period)))
}
