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
	Symbol    string          `json:"symbol" db:"symbol"`
	Exchange  string          `json:"exchange" db:"exchange"`
	Period    string          `json:"period" db:"period"` // "1m", "5m"
	Open      decimal.Decimal `json:"o" db:"open"`
	High      decimal.Decimal `json:"h" db:"high"`
	Low       decimal.Decimal `json:"l" db:"low"`
	Close     decimal.Decimal `json:"c" db:"close"`
	Volume    decimal.Decimal `json:"v" db:"volume"`
	Timestamp time.Time       `json:"t" db:"time"`
}

// OrderBook 代表深度快照 (用于回测时的高精度模拟)
type OrderBook struct {
	Symbol    string      `json:"s"`
	Timestamp time.Time   `json:"t"`
	Bids      [][2]string `json:"b"` // 使用 string 防止精度丢失，[Price, Amount]
	Asks      [][2]string `json:"a"`
}
