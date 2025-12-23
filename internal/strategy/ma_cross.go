package strategy

import (
	"market-ingestor/internal/model"
	"sync"

	"github.com/shopspring/decimal"
)

// MACrossStrategy 双均线策略
type MACrossStrategy struct {
	mu          sync.Mutex
	candles     []model.KLine
	shortPeriod int
	longPeriod  int
	lastAction  Action
}

func NewMACrossStrategy(shortPeriod, longPeriod int) *MACrossStrategy {
	return &MACrossStrategy{
		shortPeriod: shortPeriod,
		longPeriod:  longPeriod,
		candles:     make([]model.KLine, 0),
		lastAction:  ActionHold,
	}
}

func (s *MACrossStrategy) Name() string {
	return "MA_Cross"
}

func (s *MACrossStrategy) OnCandle(candle model.KLine) Action {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.candles = append(s.candles, candle)
	if len(s.candles) > s.longPeriod+1 {
		s.candles = s.candles[1:]
	}

	if len(s.candles) < s.longPeriod+1 {
		return ActionHold
	}

	shortMA := s.calculateMA(s.shortPeriod, 0)
	longMA := s.calculateMA(s.longPeriod, 0)
	prevShortMA := s.calculateMA(s.shortPeriod, 1)
	prevLongMA := s.calculateMA(s.longPeriod, 1)

	// Golden Cross
	if prevShortMA.LessThanOrEqual(prevLongMA) && shortMA.GreaterThan(longMA) {
		s.lastAction = ActionBuy
		return ActionBuy
	}
	// Death Cross
	if prevShortMA.GreaterThanOrEqual(prevLongMA) && shortMA.LessThan(longMA) {
		s.lastAction = ActionSell
		return ActionSell
	}

	return ActionHold
}

func (s *MACrossStrategy) calculateMA(period int, offset int) decimal.Decimal {
	sum := decimal.Zero
	end := len(s.candles) - offset
	start := end - period
	for i := start; i < end; i++ {
		sum = sum.Add(s.candles[i].Close)
	}
	return sum.Div(decimal.NewFromInt(int64(period)))
}
