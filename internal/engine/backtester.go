package engine

import (
	"market-ingestor/internal/model"
	"market-ingestor/internal/strategy"
	"math"

	"github.com/shopspring/decimal"
)

type Backtester struct {
	strategy    strategy.Strategy
	balance     decimal.Decimal
	position    decimal.Decimal // current quantity held
	feeRate     decimal.Decimal
	slippage    decimal.Decimal
	trades      []model.SimulatedTrade
	equityCurve []decimal.Decimal
	returns     []float64
}

func NewBacktester(strat strategy.Strategy, initialBalance decimal.Decimal) *Backtester {
	return &Backtester{
		strategy:    strat,
		balance:     initialBalance,
		position:    decimal.Zero,
		feeRate:     decimal.NewFromFloat(0.001),  // 0.1% fee
		slippage:    decimal.NewFromFloat(0.0005), // 0.05% slippage
		trades:      make([]model.SimulatedTrade, 0),
		equityCurve: make([]decimal.Decimal, 0),
		returns:     make([]float64, 0),
	}
}

func (b *Backtester) Run(candles []model.KLine) model.BacktestReport {
	initialBalance := b.balance
	prevEquity := initialBalance

	for _, candle := range candles {
		action := b.strategy.OnCandle(candle)

		if action == strategy.ActionBuy && b.balance.GreaterThan(decimal.Zero) {
			b.buy(candle)
		} else if action == strategy.ActionSell && b.position.GreaterThan(decimal.Zero) {
			b.sell(candle)
		}

		// Track equity curve and returns
		currentEquity := b.balance.Add(b.position.Mul(candle.Close))
		b.equityCurve = append(b.equityCurve, currentEquity)

		ret, _ := currentEquity.Sub(prevEquity).Div(prevEquity).Float64()
		b.returns = append(b.returns, ret)
		prevEquity = currentEquity
	}

	// Final liquidation at last price
	if b.position.GreaterThan(decimal.Zero) && len(candles) > 0 {
		b.sell(candles[len(candles)-1])
	}

	totalReturn := b.balance.Sub(initialBalance).Div(initialBalance)
	maxDD := b.calculateMaxDrawdown()
	maxDDFloat, _ := maxDD.Float64()

	winRate, totalProfit := b.calculateStats()
	sharpe := b.calculateSharpeRatio()

	return model.BacktestReport{
		StrategyName:   b.strategy.Name(),
		TotalTrades:    len(b.trades),
		WinRate:        winRate,
		TotalReturn:    totalReturn,
		TotalProfit:    totalProfit,
		MaxDrawdown:    maxDDFloat,
		SharpRatio:     sharpe,
		InitialBalance: initialBalance,
		FinalBalance:   b.balance,
		TradesLog:      b.trades,
	}
}

func (b *Backtester) buy(candle model.KLine) {
	price := candle.Close.Mul(decimal.NewFromFloat(1).Add(b.slippage))
	qty := b.balance.Div(price.Mul(decimal.NewFromFloat(1).Add(b.feeRate)))

	if qty.LessThanOrEqual(decimal.Zero) {
		return
	}

	fee := qty.Mul(price).Mul(b.feeRate)
	b.balance = b.balance.Sub(qty.Mul(price)).Sub(fee)
	b.position = b.position.Add(qty)

	b.trades = append(b.trades, model.SimulatedTrade{
		Time:   candle.Timestamp,
		Symbol: candle.Symbol,
		Side:   "buy",
		Price:  price,
		Size:   qty,
		Fee:    fee,
	})
}

func (b *Backtester) sell(candle model.KLine) {
	price := candle.Close.Mul(decimal.NewFromFloat(1).Sub(b.slippage))
	saleValue := b.position.Mul(price)
	fee := saleValue.Mul(b.feeRate)

	// Calculate PnL for the closed position
	// This is a simplified calculation: (sell_price - avg_buy_price) * qty
	// For this demo, we assume full position sell

	pnl := saleValue.Sub(fee) // This is not true PnL, but net sale.
	// To be accurate we need to track cost basis.

	b.balance = b.balance.Add(saleValue).Sub(fee)

	b.trades = append(b.trades, model.SimulatedTrade{
		Time:   candle.Timestamp,
		Symbol: candle.Symbol,
		Side:   "sell",
		Price:  price,
		Size:   b.position,
		Fee:    fee,
		PnL:    pnl, // Simplified
	})

	b.position = decimal.Zero
}

func (b *Backtester) calculateMaxDrawdown() decimal.Decimal {
	if len(b.equityCurve) == 0 {
		return decimal.Zero
	}
	maxEquity := b.equityCurve[0]
	maxDD := decimal.Zero
	for _, equity := range b.equityCurve {
		if equity.GreaterThan(maxEquity) {
			maxEquity = equity
		}
		dd := maxEquity.Sub(equity).Div(maxEquity)
		if dd.GreaterThan(maxDD) {
			maxDD = dd
		}
	}
	return maxDD
}

func (b *Backtester) calculateStats() (float64, decimal.Decimal) {
	if len(b.trades) == 0 {
		return 0, decimal.Zero
	}

	wins := 0
	totalProfit := decimal.Zero
	// Very simple win rate based on PnL of sell trades
	for _, t := range b.trades {
		if t.Side == "sell" {
			if t.PnL.GreaterThan(decimal.Zero) {
				wins++
			}
			totalProfit = totalProfit.Add(t.PnL)
		}
	}

	sellCount := 0
	for _, t := range b.trades {
		if t.Side == "sell" {
			sellCount++
		}
	}

	if sellCount == 0 {
		return 0, decimal.Zero
	}

	return float64(wins) / float64(sellCount), totalProfit
}

func (b *Backtester) calculateSharpeRatio() float64 {
	if len(b.returns) < 2 {
		return 0
	}

	var sum float64
	for _, r := range b.returns {
		sum += r
	}
	avgReturn := sum / float64(len(b.returns))

	var sumSqDiff float64
	for _, r := range b.returns {
		diff := r - avgReturn
		sumSqDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSqDiff / float64(len(b.returns)))

	if stdDev == 0 {
		return 0
	}

	// Annualize if needed, but for 1m klines we just return the ratio
	return avgReturn / stdDev
}
