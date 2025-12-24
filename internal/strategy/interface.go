package strategy

import (
	"quant-trader/internal/model"
)

type Action string

const (
	ActionBuy  Action = "buy"
	ActionSell Action = "sell"
	ActionHold Action = "hold"
)

type Strategy interface {
	Name() string
	OnCandle(candle model.KLine) Action
}
