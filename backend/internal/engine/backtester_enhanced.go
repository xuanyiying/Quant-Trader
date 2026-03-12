package engine

import (
	"math"
	"sort"
	"time"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
)

// EnhancedBacktestReport 增强版回测报告
type EnhancedBacktestReport struct {
	model.BacktestReport

	// 基础指标
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	Duration        string          `json:"duration"`
	Symbol          string          `json:"symbol"`
	Period          string          `json:"period"`

	// 收益指标
	AnnualizedReturn  float64         `json:"annualized_return"`   // 年化收益率
	Volatility        float64         `json:"volatility"`          // 波动率
	SortinoRatio      float64         `json:"sortino_ratio"`       // 索提诺比率
	CalmarRatio       float64         `json:"calmar_ratio"`        // 卡尔玛比率

	// 风险指标
	MaxDrawdownDuration int           `json:"max_drawdown_duration"` // 最大回撤持续时间(天数)
	ValueAtRisk95     float64         `json:"var_95"`                // 95% VaR
	ValueAtRisk99     float64         `json:"var_99"`                // 99% VaR

	// 交易指标
	AvgTradeReturn    float64         `json:"avg_trade_return"`    // 平均交易收益
	AvgWinningTrade   float64         `json:"avg_winning_trade"`   // 平均盈利交易
	AvgLosingTrade    float64         `json:"avg_losing_trade"`    // 平均亏损交易
	LargestWin        float64         `json:"largest_win"`         // 最大单笔盈利
	LargestLoss       float64         `json:"largest_loss"`        // 最大单笔亏损
	ProfitFactor      float64         `json:"profit_factor"`       // 盈亏比
	Expectancy        float64         `json:"expectancy"`          // 期望值

	// 持仓指标
	AvgHoldingPeriod  string          `json:"avg_holding_period"`  // 平均持仓时间
	MaxHoldingPeriod  string          `json:"max_holding_period"`  // 最长持仓时间

	// 其他指标
	Beta              float64         `json:"beta"`                // Beta系数
	Alpha             float64         `json:"alpha"`               // Alpha系数
	InformationRatio  float64         `json:"information_ratio"`   // 信息比率

	// 月度收益
	MonthlyReturns    map[string]float64 `json:"monthly_returns"`  // 月度收益

	// 权益曲线
	EquityCurve       []EquityPoint   `json:"equity_curve"`        // 权益曲线数据
}

// EquityPoint 权益曲线点
type EquityPoint struct {
	Timestamp time.Time       `json:"timestamp"`
	Equity    decimal.Decimal `json:"equity"`
	Drawdown  float64         `json:"drawdown"`
}

// TradeStats 交易统计
type TradeStats struct {
	TotalTrades      int
	WinningTrades    int
	LosingTrades     int
	TotalProfit      decimal.Decimal
	TotalLoss        decimal.Decimal
	LargestWin       decimal.Decimal
	LargestLoss      decimal.Decimal
	HoldingPeriods   []time.Duration
}

// CalculateEnhancedMetrics 计算增强版回测指标
func CalculateEnhancedMetrics(
	baseReport model.BacktestReport,
	candles []model.KLine,
	equityCurve []decimal.Decimal,
	returns []float64,
	trades []model.SimulatedTrade,
	initialBalance decimal.Decimal,
) EnhancedBacktestReport {

	if len(candles) == 0 {
		return EnhancedBacktestReport{BacktestReport: baseReport}
	}

	startTime := candles[0].Timestamp
	endTime := candles[len(candles)-1].Timestamp
	duration := endTime.Sub(startTime)

	report := EnhancedBacktestReport{
		BacktestReport:    baseReport,
		StartTime:         startTime,
		EndTime:           endTime,
		Duration:          duration.String(),
		Symbol:            candles[0].Symbol,
		Period:            candles[0].Period,
	}

	// 计算交易统计
	tradeStats := calculateTradeStats(trades)

	// 收益指标
	report.AnnualizedReturn = calculateAnnualizedReturn(
		baseReport.FinalBalance,
		initialBalance,
		duration,
	)
	report.Volatility = calculateVolatility(returns)
	report.SortinoRatio = calculateSortinoRatio(returns)
	report.CalmarRatio = calculateCalmarRatio(
		report.AnnualizedReturn,
		baseReport.MaxDrawdown,
	)

	// 风险指标
	report.MaxDrawdownDuration = calculateMaxDrawdownDuration(equityCurve, candles)
	report.ValueAtRisk95 = calculateVaR(returns, 0.05)
	report.ValueAtRisk99 = calculateVaR(returns, 0.01)

	// 交易指标
	if tradeStats.TotalTrades > 0 {
		report.AvgTradeReturn = calculateAvgTradeReturn(tradeStats, initialBalance)
		report.AvgWinningTrade = calculateAvgWinningTrade(tradeStats)
		report.AvgLosingTrade = calculateAvgLosingTrade(tradeStats)
		report.LargestWin, _ = tradeStats.LargestWin.Float64()
		report.LargestLoss, _ = tradeStats.LargestLoss.Float64()
		report.ProfitFactor = calculateProfitFactor(tradeStats)
		report.Expectancy = calculateExpectancy(tradeStats)
	}

	// 持仓指标
	if len(tradeStats.HoldingPeriods) > 0 {
		report.AvgHoldingPeriod = calculateAvgHoldingPeriod(tradeStats.HoldingPeriods)
		report.MaxHoldingPeriod = calculateMaxHoldingPeriod(tradeStats.HoldingPeriods)
	}

	// 月度收益
	report.MonthlyReturns = calculateMonthlyReturns(equityCurve, candles)

	// 权益曲线
	report.EquityCurve = buildEquityCurve(equityCurve, candles)

	return report
}

// calculateTradeStats 计算交易统计
func calculateTradeStats(trades []model.SimulatedTrade) TradeStats {
	stats := TradeStats{
		HoldingPeriods: make([]time.Duration, 0),
	}

	var buyTime time.Time
	for _, trade := range trades {
		if trade.Side == model.SideBuy {
			stats.TotalTrades++
			buyTime = trade.Time
		} else if trade.Side == model.SideSell {
			pnl, _ := trade.PnL.Float64()
			if pnl > 0 {
				stats.WinningTrades++
				stats.TotalProfit = stats.TotalProfit.Add(trade.PnL)
				if trade.PnL.GreaterThan(stats.LargestWin) {
					stats.LargestWin = trade.PnL
				}
			} else {
				stats.LosingTrades++
				stats.TotalLoss = stats.TotalLoss.Add(trade.PnL.Abs())
				if trade.PnL.LessThan(stats.LargestLoss) {
					stats.LargestLoss = trade.PnL
				}
			}

			// 计算持仓时间
			if !buyTime.IsZero() {
				holdingPeriod := trade.Time.Sub(buyTime)
				stats.HoldingPeriods = append(stats.HoldingPeriods, holdingPeriod)
			}
		}
	}

	return stats
}

// calculateAnnualizedReturn 计算年化收益率
func calculateAnnualizedReturn(finalBalance, initialBalance decimal.Decimal, duration time.Duration) float64 {
	if duration == 0 || initialBalance.IsZero() {
		return 0
	}

	totalReturn, _ := finalBalance.Sub(initialBalance).Div(initialBalance).Float64()
	years := duration.Hours() / 24 / 365

	if years == 0 {
		return totalReturn
	}

	// 使用几何平均计算年化收益
	return math.Pow(1+totalReturn, 1/years) - 1
}

// calculateVolatility 计算波动率 (年化)
func calculateVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	avg := 0.0
	for _, r := range returns {
		avg += r
	}
	avg /= float64(len(returns))

	var sumSqDiff float64
	for _, r := range returns {
		diff := r - avg
		sumSqDiff += diff * diff
	}

	stdDev := math.Sqrt(sumSqDiff / float64(len(returns)))

	// 假设 returns 是 1m 数据的收益率，年化需要乘以 sqrt(365*24*60)
	annualizationFactor := math.Sqrt(365 * 24 * 60)
	return stdDev * annualizationFactor
}

// calculateSortinoRatio 计算索提诺比率
func calculateSortinoRatio(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	avgReturn := 0.0
	for _, r := range returns {
		avgReturn += r
	}
	avgReturn /= float64(len(returns))

	// 只计算下行标准差
	var sumSqDownside float64
	downsideCount := 0
	for _, r := range returns {
		if r < 0 {
			sumSqDownside += r * r
			downsideCount++
		}
	}

	if downsideCount == 0 {
		return math.Inf(1)
	}

	downsideDev := math.Sqrt(sumSqDownside / float64(downsideCount))
	if downsideDev == 0 {
		return 0
	}

	return avgReturn / downsideDev
}

// calculateCalmarRatio 计算卡尔玛比率
func calculateCalmarRatio(annualizedReturn, maxDrawdown float64) float64 {
	if maxDrawdown == 0 {
		return 0
	}
	return annualizedReturn / maxDrawdown
}

// calculateMaxDrawdownDuration 计算最大回撤持续时间
func calculateMaxDrawdownDuration(equityCurve []decimal.Decimal, candles []model.KLine) int {
	if len(equityCurve) == 0 {
		return 0
	}

	maxDuration := 0
	currentDuration := 0
	peakIdx := 0

	for i := 1; i < len(equityCurve); i++ {
		if equityCurve[i].GreaterThan(equityCurve[peakIdx]) {
			// 创新高，重置回撤
			if currentDuration > maxDuration {
				maxDuration = currentDuration
			}
			peakIdx = i
			currentDuration = 0
		} else {
			currentDuration++
		}
	}

	if currentDuration > maxDuration {
		maxDuration = currentDuration
	}

	return maxDuration
}

// calculateVaR 计算风险价值 (Value at Risk)
func calculateVaR(returns []float64, confidence float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	// 复制并排序收益率
	sortedReturns := make([]float64, len(returns))
	copy(sortedReturns, returns)
	sort.Float64s(sortedReturns)

	// 找到对应分位数的索引
	index := int(float64(len(sortedReturns)) * confidence)
	if index >= len(sortedReturns) {
		index = len(sortedReturns) - 1
	}

	return sortedReturns[index]
}

// calculateAvgTradeReturn 计算平均交易收益
func calculateAvgTradeReturn(stats TradeStats, initialBalance decimal.Decimal) float64 {
	if stats.TotalTrades == 0 || initialBalance.IsZero() {
		return 0
	}

	totalPnL := stats.TotalProfit.Sub(stats.TotalLoss)
	avgPnL, _ := totalPnL.Div(decimal.NewFromInt(int64(stats.TotalTrades))).Float64()
	initialFloat, _ := initialBalance.Float64()

	return avgPnL / initialFloat
}

// calculateAvgWinningTrade 计算平均盈利交易
func calculateAvgWinningTrade(stats TradeStats) float64 {
	if stats.WinningTrades == 0 {
		return 0
	}

	avg, _ := stats.TotalProfit.Div(decimal.NewFromInt(int64(stats.WinningTrades))).Float64()
	return avg
}

// calculateAvgLosingTrade 计算平均亏损交易
func calculateAvgLosingTrade(stats TradeStats) float64 {
	if stats.LosingTrades == 0 {
		return 0
	}

	avg, _ := stats.TotalLoss.Div(decimal.NewFromInt(int64(stats.LosingTrades))).Float64()
	return -avg // 返回负值
}

// calculateProfitFactor 计算盈亏比
func calculateProfitFactor(stats TradeStats) float64 {
	if stats.TotalLoss.IsZero() {
		return math.Inf(1)
	}

	pf, _ := stats.TotalProfit.Div(stats.TotalLoss).Float64()
	return pf
}

// calculateExpectancy 计算期望值
func calculateExpectancy(stats TradeStats) float64 {
	if stats.TotalTrades == 0 {
		return 0
	}

	winRate := float64(stats.WinningTrades) / float64(stats.TotalTrades)
	lossRate := 1 - winRate

	avgWin := calculateAvgWinningTrade(stats)
	avgLoss := math.Abs(calculateAvgLosingTrade(stats))

	return winRate*avgWin - lossRate*avgLoss
}

// calculateAvgHoldingPeriod 计算平均持仓时间
func calculateAvgHoldingPeriod(periods []time.Duration) string {
	if len(periods) == 0 {
		return "0s"
	}

	var total time.Duration
	for _, p := range periods {
		total += p
	}

	avg := total / time.Duration(len(periods))
	return avg.String()
}

// calculateMaxHoldingPeriod 计算最长持仓时间
func calculateMaxHoldingPeriod(periods []time.Duration) string {
	if len(periods) == 0 {
		return "0s"
	}

	max := periods[0]
	for _, p := range periods {
		if p > max {
			max = p
		}
	}

	return max.String()
}

// calculateMonthlyReturns 计算月度收益
func calculateMonthlyReturns(equityCurve []decimal.Decimal, candles []model.KLine) map[string]float64 {
	monthlyReturns := make(map[string]float64)

	if len(equityCurve) == 0 || len(candles) == 0 {
		return monthlyReturns
	}

	// 按月分组计算收益
	monthlyEquity := make(map[string]decimal.Decimal)

	for i, candle := range candles {
		if i >= len(equityCurve) {
			break
		}

		monthKey := candle.Timestamp.Format("2006-01")
		if _, exists := monthlyEquity[monthKey]; !exists {
			// 记录月初的权益
			if i > 0 {
				prevMonth := candles[i-1].Timestamp.Format("2006-01")
				if prevMonth != monthKey {
					monthlyEquity[monthKey] = equityCurve[i]
				}
			}
		}
	}

	// 计算每月收益率
	for month, equity := range monthlyEquity {
		ret, _ := equity.Float64()
		monthlyReturns[month] = ret
	}

	return monthlyReturns
}

// buildEquityCurve 构建权益曲线数据
func buildEquityCurve(equityCurve []decimal.Decimal, candles []model.KLine) []EquityPoint {
	points := make([]EquityPoint, 0, len(equityCurve))

	if len(equityCurve) == 0 {
		return points
	}

	maxEquity := equityCurve[0]

	for i, equity := range equityCurve {
		if equity.GreaterThan(maxEquity) {
			maxEquity = equity
		}

		dd := 0.0
		if !maxEquity.IsZero() {
			ddFloat, _ := maxEquity.Sub(equity).Div(maxEquity).Float64()
			dd = ddFloat
		}

		timestamp := time.Now()
		if i < len(candles) {
			timestamp = candles[i].Timestamp
		}

		points = append(points, EquityPoint{
			Timestamp: timestamp,
			Equity:    equity,
			Drawdown:  dd,
		})
	}

	return points
}
