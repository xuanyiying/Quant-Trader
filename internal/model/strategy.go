package model

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// Strategy 策略配置实体
type Strategy struct {
	ID        int64           `json:"id" db:"id"`
	UserID    int64           `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Type      string          `json:"type" db:"type"`
	Config    json.RawMessage `json:"config" db:"config"` // 灵活存储配置
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// BacktestReport 回测结果报告
type BacktestReport struct {
	StrategyName   string          `json:"strategy_name"`
	TotalTrades    int             `json:"total_trades"`
	WinRate        float64         `json:"win_rate"`
	TotalReturn    decimal.Decimal `json:"total_return"`
	TotalProfit    decimal.Decimal `json:"total_profit"` // 净利润
	MaxDrawdown    float64         `json:"max_drawdown"` // 最大回撤
	SharpRatio     float64         `json:"sharp_ratio"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
	FinalBalance   decimal.Decimal `json:"final_balance"`
	TradesLog      []SimulatedTrade `json:"trades_log"`   // 交易明细
}

// SimulatedTrade 回测中的单笔交易记录
type SimulatedTrade struct {
	Time   time.Time       `json:"time"`
	Symbol string          `json:"symbol"`
	Side   string          `json:"side"` // "buy", "sell"
	Price  decimal.Decimal `json:"price"`
	Size   decimal.Decimal `json:"size"`
	Fee    decimal.Decimal `json:"fee"`
	PnL    decimal.Decimal `json:"pnl"` // Profit and Loss
}
