package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BacktestRun struct {
	ID          int64           `gorm:"primaryKey" json:"id"`
	StrategyID  int64           `gorm:"index;not null" json:"strategy_id"`
	UserID      int64           `gorm:"index;not null" json:"user_id"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     time.Time       `json:"end_time"`
	InitialCap  decimal.Decimal `gorm:"type:decimal(20,8)" json:"initial_cap"`
	FinalCap    decimal.Decimal `gorm:"type:decimal(20,8)" json:"final_cap"`
	TotalReturn decimal.Decimal `gorm:"type:decimal(10,4)" json:"total_return"`
	SharpeRatio decimal.Decimal `gorm:"type:decimal(10,4)" json:"sharpe_ratio"`
	MaxDrawdown decimal.Decimal `gorm:"type:decimal(10,4)" json:"max_drawdown"`
	WinRate     decimal.Decimal `gorm:"type:decimal(10,4)" json:"win_rate"`
	Trades      int             `json:"trades"`
	Results     string          `gorm:"type:jsonb" json:"results"`
	CreatedAt   time.Time       `json:"created_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (BacktestRun) TableName() string {
	return "backtest_runs"
}
