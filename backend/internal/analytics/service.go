package analytics

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrNoDataAvailable   = errors.New("no data available for calculation")
	ErrInvalidIterations = errors.New("invalid iteration count")
	ErrInvalidDays       = errors.New("invalid days count")
)

type AnalyticsService struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAnalyticsService(db *pgxpool.Pool, logger *zap.Logger) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		logger: logger,
	}
}

type PerformanceReport struct {
	TotalReturn      decimal.Decimal `json:"total_return"`
	MaxDrawdown      decimal.Decimal `json:"max_drawdown"`
	SharpeRatio      float64         `json:"sharpe_ratio"`
	WinRate          float64         `json:"win_rate"`
	TotalTrades      int64           `json:"total_trades"`
	ProfitableTrades int64           `json:"profitable_trades"`
}

func (s *AnalyticsService) GetPortfolioReport(ctx context.Context, userID int64) (*PerformanceReport, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	s.logger.Info("generating portfolio report", zap.Int64("user_id", userID))

	var balance decimal.Decimal
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(balance, 0) 
		FROM paper_accounts 
		WHERE user_id = $1`, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("paper account not found", zap.Int64("user_id", userID))
			return &PerformanceReport{
				TotalReturn: decimal.Zero,
				MaxDrawdown: decimal.Zero,
				SharpeRatio: 0,
				WinRate:     0,
			}, nil
		}
		s.logger.Error("failed to fetch account balance", zap.Error(err), zap.Int64("user_id", userID))
		return nil, err
	}

	rows, err := s.db.Query(ctx, `
		SELECT symbol, qty, avg_price 
		FROM paper_positions 
		WHERE user_id = $1 AND qty > 0`, userID)
	if err != nil {
		s.logger.Error("failed to fetch positions", zap.Error(err), zap.Int64("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	totalMarketValue := balance
	var positions []struct {
		Symbol   string
		Quantity decimal.Decimal
		AvgPrice decimal.Decimal
	}

	for rows.Next() {
		var pos struct {
			Symbol   string
			Quantity decimal.Decimal
			AvgPrice decimal.Decimal
		}
		if err := rows.Scan(&pos.Symbol, &pos.Quantity, &pos.AvgPrice); err != nil {
			s.logger.Warn("failed to scan position", zap.Error(err))
			continue
		}
		positions = append(positions, pos)
	}

	for _, pos := range positions {
		currentPrice, err := s.getCurrentPrice(ctx, pos.Symbol)
		if err != nil {
			s.logger.Warn("failed to get current price, using average price",
				zap.String("symbol", pos.Symbol), zap.Error(err))
			currentPrice = pos.AvgPrice
		}
		totalMarketValue = totalMarketValue.Add(pos.Quantity.Mul(currentPrice))
	}

	initialInvestment, err := s.getInitialInvestment(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to get initial investment, using current market value as baseline")
		initialInvestment = totalMarketValue
	}

	totalReturn := decimal.Zero
	if initialInvestment.GreaterThan(decimal.Zero) {
		totalReturn = totalMarketValue.Sub(initialInvestment).Div(initialInvestment).Mul(decimal.NewFromInt(100))
	}

	maxDrawdown, err := s.calculateMaxDrawdown(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to calculate max drawdown", zap.Error(err))
		maxDrawdown = decimal.Zero
	}

	sharpeRatio, err := s.calculateSharpeRatio(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to calculate sharpe ratio", zap.Error(err))
		sharpeRatio = 0
	}

	winRate, totalTrades, profitableTrades, err := s.calculateWinRate(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to calculate win rate", zap.Error(err))
		winRate = 0
		totalTrades = 0
		profitableTrades = 0
	}

	s.logger.Info("portfolio report generated",
		zap.Int64("user_id", userID),
		zap.String("total_return", totalReturn.String()),
		zap.Float64("sharpe_ratio", sharpeRatio),
		zap.Float64("win_rate", winRate))

	return &PerformanceReport{
		TotalReturn:      totalReturn,
		MaxDrawdown:      maxDrawdown,
		SharpeRatio:      sharpeRatio,
		WinRate:          winRate,
		TotalTrades:      totalTrades,
		ProfitableTrades: profitableTrades,
	}, nil
}

func (s *AnalyticsService) getCurrentPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	var price decimal.Decimal
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(close, 0) 
		FROM klines 
		WHERE symbol = $1 
		ORDER BY time DESC 
		LIMIT 1`, symbol).Scan(&price)
	if err != nil {
		return decimal.Zero, err
	}
	return price, nil
}

func (s *AnalyticsService) getInitialInvestment(ctx context.Context, userID int64) (decimal.Decimal, error) {
	var initialBalance decimal.Decimal
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(initial_balance, 0) 
		FROM paper_accounts 
		WHERE user_id = $1`, userID).Scan(&initialBalance)
	if err != nil {
		return decimal.Zero, err
	}
	return initialBalance, nil
}

func (s *AnalyticsService) calculateMaxDrawdown(ctx context.Context, userID int64) (decimal.Decimal, error) {
	rows, err := s.db.Query(ctx, `
		SELECT time, balance 
		FROM paper_account_history 
		WHERE user_id = $1 
		ORDER BY time ASC`, userID)
	if err != nil {
		return decimal.Zero, err
	}
	defer rows.Close()

	if !rows.Next() {
		return decimal.Zero, nil
	}

	var peak decimal.Decimal
	var maxDrawdown decimal.Decimal

	var timestamp time.Time
	var balance decimal.Decimal
	if err := rows.Scan(&timestamp, &balance); err != nil {
		return decimal.Zero, err
	}
	peak = balance

	for rows.Next() {
		if err := rows.Scan(&timestamp, &balance); err != nil {
			continue
		}
		if balance.GreaterThan(peak) {
			peak = balance
		}
		drawdown := peak.Sub(balance).Div(peak).Mul(decimal.NewFromInt(100))
		if drawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown, nil
}

func (s *AnalyticsService) calculateSharpeRatio(ctx context.Context, userID int64) (float64, error) {
	rows, err := s.db.Query(ctx, `
		SELECT time, balance 
		FROM paper_account_history 
		WHERE user_id = $1 
		ORDER BY time ASC`, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var previousBalance decimal.Decimal
	var returns []float64

	for rows.Next() {
		var timestamp time.Time
		var balance decimal.Decimal
		if err := rows.Scan(&timestamp, &balance); err != nil {
			continue
		}
		if previousBalance.GreaterThan(decimal.Zero) {
			dailyReturn, ok := balance.Sub(previousBalance).Div(previousBalance).Float64()
			if ok {
				returns = append(returns, dailyReturn)
			}
		}
		previousBalance = balance
	}

	if len(returns) == 0 {
		return 0, nil
	}

	riskFreeRate := 0.04 / 365
	return CalculateSharpeRatio(returns, riskFreeRate), nil
}

func (s *AnalyticsService) calculateWinRate(ctx context.Context, userID int64) (float64, int64, int64, error) {
	var totalTrades int64
	var profitableTrades int64

	err := s.db.QueryRow(ctx, `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE pnl > 0) as profitable
		FROM paper_trades 
		WHERE user_id = $1`, userID).Scan(&totalTrades, &profitableTrades)
	if err != nil {
		return 0, 0, 0, err
	}

	if totalTrades == 0 {
		return 0, totalTrades, profitableTrades, nil
	}

	winRate := float64(profitableTrades) / float64(totalTrades)
	return winRate, totalTrades, profitableTrades, nil
}

func (s *AnalyticsService) MonteCarloSimulation(ctx context.Context, returns []float64, iterations int, days int) ([][]float64, error) {
	if len(returns) == 0 {
		return nil, ErrNoDataAvailable
	}
	if iterations <= 0 || iterations > 10000 {
		return nil, ErrInvalidIterations
	}
	if days <= 0 || days > 3650 {
		return nil, ErrInvalidDays
	}

	s.logger.Info("running Monte Carlo simulation",
		zap.Int("iterations", iterations),
		zap.Int("days", days),
		zap.Int("historical_returns", len(returns)))

	results := make([][]float64, iterations)
	for i := 0; i < iterations; i++ {
		path := make([]float64, days)
		price := 1.0
		for d := 0; d < days; d++ {
			r := returns[int(time.Now().UnixNano())%len(returns)]
			price = price * (1 + r)
			path[d] = price
		}
		results[i] = path
	}

	s.logger.Info("Monte Carlo simulation completed", zap.Int("iterations", iterations))
	return results, nil
}

func CalculateSharpeRatio(returns []float64, riskFreeRate float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	sum := 0.0
	for _, r := range returns {
		sum += r
	}
	mean := sum / float64(len(returns))

	variance := 0.0
	for _, r := range returns {
		diff := r - mean
		variance += diff * diff
	}
	stdDev := math.Sqrt(variance / float64(len(returns)))

	if stdDev == 0 {
		return 0
	}

	return (mean - riskFreeRate) / stdDev
}
