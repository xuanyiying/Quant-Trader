package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Trade 代表一笔实时成交
type Trade struct {
	ID        string          `json:"id" db:"trade_id"`
	Symbol    string          `json:"symbol" db:"symbol"`
	Exchange  string          `json:"exchange" db:"exchange"`
	Price     decimal.Decimal `json:"price" db:"price"`
	Amount    decimal.Decimal `json:"amount" db:"amount"`
	Side      string          `json:"side" db:"side"` // "buy" or "sell"
	Timestamp time.Time       `json:"ts" db:"time"`
}

// KLine (Candle) 代表一根K线
type KLine struct {
	ID        int64           `gorm:"primaryKey" json:"id"`
	Symbol    string          `gorm:"index:idx_symbol_period,priority:1" json:"symbol"`
	Exchange  string          `gorm:"index:idx_symbol_period,priority:2" json:"exchange"`
	Period    string          `gorm:"index:idx_symbol_period,priority:3" json:"period"`
	Open      decimal.Decimal `gorm:"type:decimal(20,8)" json:"o"`
	High      decimal.Decimal `gorm:"type:decimal(20,8)" json:"h"`
	Low       decimal.Decimal `gorm:"type:decimal(20,8)" json:"l"`
	Close     decimal.Decimal `gorm:"type:decimal(20,8)" json:"c"`
	Volume    decimal.Decimal `gorm:"type:decimal(20,8)" json:"v"`
	Timestamp time.Time      `gorm:"index" json:"t"`
}

func (KLine) TableName() string {
	return "klines"
}

// SupportedPeriods 定义系统支持的 K 线周期
var SupportedPeriods = []string{Period1m, Period5m, Period15m, Period1h, Period4h, Period1d}

// PeriodToDuration 将字符串周期转换为 time.Duration
func PeriodToDuration(period string) time.Duration {
	switch period {
	case Period1m:
		return time.Minute
	case Period5m:
		return 5 * time.Minute
	case Period15m:
		return 15 * time.Minute
	case Period1h:
		return time.Hour
	case Period4h:
		return 4 * time.Hour
	case Period1d:
		return 24 * time.Hour
	default:
		return time.Minute
	}
}

// OrderBook 代表深度快照 (用于回测时的高精度模拟)
type OrderBook struct {
	Symbol    string      `json:"s"`
	Timestamp time.Time   `json:"t"`
	Bids      [][2]string `json:"b"` // 使用 string 防止精度丢失，[Price, Amount]
	Asks      [][2]string `json:"a"`
}
