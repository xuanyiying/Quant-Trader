package model

import (
	"github.com/shopspring/decimal"
)

// Order represents a unified order structure for both matching and paper trading
type Order struct {
	ID           string          `json:"id"`
	UserID       int64           `json:"user_id,omitempty"`
	AccountID    string          `json:"account_id,omitempty"`
	Symbol       string          `json:"symbol"`
	Side         string          `json:"side"` // buy, sell
	Type         string          `json:"type"` // market, limit, etc.
	Price        decimal.Decimal `json:"price"`
	Qty          decimal.Decimal `json:"qty"`
	FilledQty    decimal.Decimal `json:"filled_qty"`
	FilledPrice  decimal.Decimal `json:"filled_price"`
	Status       string          `json:"status"`
	Timestamp    int64           `json:"timestamp"`
	TriggerPrice decimal.Decimal `json:"trigger_price,omitempty"`
}

// Execution represents a successful match or execution
type Execution struct {
	ID         string          `json:"id"`
	Symbol     string          `json:"symbol"`
	Price      decimal.Decimal `json:"price"`
	Qty        decimal.Decimal `json:"qty"`
	MakerOrder string          `json:"maker_order_id,omitempty"`
	TakerOrder string          `json:"taker_order_id,omitempty"`
	Side       string          `json:"side"` // taker's side
	Timestamp  int64           `json:"timestamp"`
}

// Event is the unified event format for frontend consumption
type Event struct {
	Seq       uint64     `json:"seq"`
	Type      string     `json:"type"` // use EventType constants
	Timestamp int64      `json:"timestamp"`
	Symbol    string     `json:"symbol"`
	Order     *Order     `json:"order,omitempty"`
	Execution *Execution `json:"execution,omitempty"`
	Reason    string     `json:"reason,omitempty"`
}
