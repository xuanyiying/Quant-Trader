package matching

import "github.com/shopspring/decimal"

type EngineSnapshot struct {
	Symbol         string          `json:"symbol"`
	CapturedAt     int64           `json:"captured_at"`
	Seq            uint64          `json:"seq"`
	State          MarketState     `json:"state"`
	LastPrice      decimal.Decimal `json:"last_price"`
	LastPriceTicks int64           `json:"last_price_ticks"`
	Config         SymbolConfig    `json:"config"`
	OpenOrders     []OrderRecord   `json:"open_orders"`
	Triggers       []OrderRecord   `json:"triggers"`
}
