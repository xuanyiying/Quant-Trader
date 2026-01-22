package matching

import (
	"time"

	"github.com/shopspring/decimal"
)

// OrderSide defines the side of the order (Buy/Sell)
type OrderSide int

const (
	SideBuy OrderSide = iota
	SideSell
)

func (s OrderSide) String() string {
	if s == SideBuy {
		return "buy"
	}
	return "sell"
}

// OrderType defines the type of order
type OrderType int

const (
	TypeLimit OrderType = iota
	TypeMarket
	TypeFOK // Fill or Kill
	TypeIOC // Immediate or Cancel
	TypeStop
	TypeTakeProfit
	TypeIceberg
)

// OrderStatus represents the current state of an order
type OrderStatus int

const (
	StatusNew OrderStatus = iota
	StatusOpen
	StatusPartiallyFilled
	StatusPartiallyFilledCanceled
	StatusFilled
	StatusCanceled
	StatusRejected
)

// Order represents an instruction to buy or sell
type Order struct {
	ID                  string
	AccountID           string
	Symbol              string
	Side                OrderSide
	Type                OrderType
	Price               decimal.Decimal
	PriceTicks          int64
	Amount              decimal.Decimal // Original amount
	AmountLots          int64
	Filled              decimal.Decimal // Amount filled so far
	FilledLots          int64
	Remaining           decimal.Decimal // Amount remaining (Amount - Filled)
	RemainingLots       int64
	TriggerPrice        decimal.Decimal
	TriggerTicks        int64
	PeakSize            decimal.Decimal
	PeakLots            int64
	HiddenRemaining     decimal.Decimal
	HiddenRemainingLots int64
	Timestamp           int64 // Unix nanoseconds
	Status              OrderStatus
}

type SymbolConfig struct {
	TickSize           decimal.Decimal
	LotSize            decimal.Decimal
	MinQty             decimal.Decimal
	MinNotional        decimal.Decimal
	LimitUp            decimal.Decimal
	LimitDown          decimal.Decimal
	LargeOrderNotional decimal.Decimal
	MaxOrdersPerSecond int
}

type MarketState int

const (
	StatePreOpen MarketState = iota
	StateOpenAuction
	StateContinuous
	StateClosing
	StateClosed
)

type RiskManager interface {
	ValidateOrder(state MarketState, lastPrice decimal.Decimal, snap OrderBookSnapshot, order *Order) error
	OnTrade(trade Trade)
}

type NoopRiskManager struct{}

func (NoopRiskManager) ValidateOrder(state MarketState, lastPrice decimal.Decimal, snap OrderBookSnapshot, order *Order) error {
	return nil
}

func (NoopRiskManager) OnTrade(trade Trade) {}

type FeeTier struct {
	Notional decimal.Decimal
	Rate     decimal.Decimal
}

type TieredFeeCalculator struct {
	Tiers       []FeeTier
	DefaultRate decimal.Decimal
}

func (c TieredFeeCalculator) Fee(notional decimal.Decimal) decimal.Decimal {
	rate := c.DefaultRate
	for i := range c.Tiers {
		if notional.GreaterThanOrEqual(c.Tiers[i].Notional) {
			rate = c.Tiers[i].Rate
		}
	}
	return notional.Mul(rate)
}

func VWAP(trades []Trade) (decimal.Decimal, decimal.Decimal) {
	totalNotional := decimal.Zero
	totalQty := decimal.Zero
	for i := range trades {
		totalNotional = totalNotional.Add(trades[i].Price.Mul(trades[i].Amount))
		totalQty = totalQty.Add(trades[i].Amount)
	}
	if totalQty.IsZero() {
		return decimal.Zero, decimal.Zero
	}
	return totalNotional.Div(totalQty), totalQty
}

type Position struct {
	Qty         decimal.Decimal
	AvgPrice    decimal.Decimal
	RealizedPnL decimal.Decimal
	FeePaid     decimal.Decimal
}

func (p *Position) ApplyTrade(side OrderSide, price, qty, fee decimal.Decimal) {
	if qty.LessThanOrEqual(decimal.Zero) {
		return
	}

	p.FeePaid = p.FeePaid.Add(fee)

	if side == SideBuy {
		if p.Qty.GreaterThanOrEqual(decimal.Zero) {
			newQty := p.Qty.Add(qty)
			if newQty.IsZero() {
				p.Qty = decimal.Zero
				p.AvgPrice = decimal.Zero
				return
			}
			cost := p.AvgPrice.Mul(p.Qty).Add(price.Mul(qty))
			p.Qty = newQty
			p.AvgPrice = cost.Div(newQty)
			return
		}
		closeQty := decimal.Min(qty, p.Qty.Abs())
		p.RealizedPnL = p.RealizedPnL.Add(p.AvgPrice.Sub(price).Mul(closeQty))
		p.Qty = p.Qty.Add(closeQty)
		openQty := qty.Sub(closeQty)
		if openQty.GreaterThan(decimal.Zero) {
			p.Qty = p.Qty.Add(openQty)
			p.AvgPrice = price
		}
		return
	}

	if p.Qty.LessThanOrEqual(decimal.Zero) {
		newQty := p.Qty.Sub(qty)
		if newQty.IsZero() {
			p.Qty = decimal.Zero
			p.AvgPrice = decimal.Zero
			return
		}
		proceeds := p.AvgPrice.Mul(p.Qty.Abs()).Add(price.Mul(qty))
		p.Qty = newQty
		p.AvgPrice = proceeds.Div(newQty.Abs())
		return
	}

	closeQty := decimal.Min(qty, p.Qty)
	p.RealizedPnL = p.RealizedPnL.Add(price.Sub(p.AvgPrice).Mul(closeQty))
	p.Qty = p.Qty.Sub(closeQty)
	openQty := qty.Sub(closeQty)
	if openQty.GreaterThan(decimal.Zero) {
		p.Qty = p.Qty.Sub(openQty)
		p.AvgPrice = price
	}
}

func (p *Position) UnrealizedPnL(mark decimal.Decimal) decimal.Decimal {
	if p.Qty.IsZero() {
		return decimal.Zero
	}
	if p.Qty.GreaterThan(decimal.Zero) {
		return mark.Sub(p.AvgPrice).Mul(p.Qty)
	}
	return p.AvgPrice.Sub(mark).Mul(p.Qty.Abs())
}

type EventType int

const (
	EventSubmit EventType = iota
	EventCancel
	EventTrigger
	EventTrade
	EventOrderUpdate
	EventReject
	EventAlert
)

type OrderRecord struct {
	ID                  string          `json:"id"`
	AccountID           string          `json:"account_id"`
	Symbol              string          `json:"symbol"`
	Side                OrderSide       `json:"side"`
	Type                OrderType       `json:"type"`
	Price               decimal.Decimal `json:"price"`
	PriceTicks          int64           `json:"price_ticks"`
	Amount              decimal.Decimal `json:"amount"`
	AmountLots          int64           `json:"amount_lots"`
	Filled              decimal.Decimal `json:"filled"`
	FilledLots          int64           `json:"filled_lots"`
	Remaining           decimal.Decimal `json:"remaining"`
	RemainingLots       int64           `json:"remaining_lots"`
	TriggerPrice        decimal.Decimal `json:"trigger_price"`
	TriggerTicks        int64           `json:"trigger_ticks"`
	PeakSize            decimal.Decimal `json:"peak_size"`
	PeakLots            int64           `json:"peak_lots"`
	HiddenRemaining     decimal.Decimal `json:"hidden_remaining"`
	HiddenRemainingLots int64           `json:"hidden_remaining_lots"`
	Timestamp           int64           `json:"timestamp"`
	Status              OrderStatus     `json:"status"`
}

func (o *Order) Record() OrderRecord {
	return OrderRecord{
		ID:                  o.ID,
		AccountID:           o.AccountID,
		Symbol:              o.Symbol,
		Side:                o.Side,
		Type:                o.Type,
		Price:               o.Price,
		PriceTicks:          o.PriceTicks,
		Amount:              o.Amount,
		AmountLots:          o.AmountLots,
		Filled:              o.Filled,
		FilledLots:          o.FilledLots,
		Remaining:           o.Remaining,
		RemainingLots:       o.RemainingLots,
		TriggerPrice:        o.TriggerPrice,
		TriggerTicks:        o.TriggerTicks,
		PeakSize:            o.PeakSize,
		PeakLots:            o.PeakLots,
		HiddenRemaining:     o.HiddenRemaining,
		HiddenRemainingLots: o.HiddenRemainingLots,
		Timestamp:           o.Timestamp,
		Status:              o.Status,
	}
}

func (r OrderRecord) ToOrder() *Order {
	o := &Order{
		ID:                  r.ID,
		AccountID:           r.AccountID,
		Symbol:              r.Symbol,
		Side:                r.Side,
		Type:                r.Type,
		Price:               r.Price,
		PriceTicks:          r.PriceTicks,
		Amount:              r.Amount,
		AmountLots:          r.AmountLots,
		Filled:              r.Filled,
		FilledLots:          r.FilledLots,
		Remaining:           r.Remaining,
		RemainingLots:       r.RemainingLots,
		TriggerPrice:        r.TriggerPrice,
		TriggerTicks:        r.TriggerTicks,
		PeakSize:            r.PeakSize,
		PeakLots:            r.PeakLots,
		HiddenRemaining:     r.HiddenRemaining,
		HiddenRemainingLots: r.HiddenRemainingLots,
		Timestamp:           r.Timestamp,
		Status:              r.Status,
	}
	if o.Remaining.IsZero() && o.Amount.GreaterThan(decimal.Zero) && o.Filled.IsZero() {
		o.Remaining = o.Amount
	}
	return o
}

type Event struct {
	Seq       uint64       `json:"seq"`
	Type      EventType    `json:"type"`
	Timestamp int64        `json:"timestamp"`
	Symbol    string       `json:"symbol"`
	Order     *OrderRecord `json:"order,omitempty"`
	OrderID   string       `json:"order_id,omitempty"`
	Trade     *Trade       `json:"trade,omitempty"`
	Reason    string       `json:"reason,omitempty"`
}

// NewOrder creates a new order instance
func NewOrder(id string, symbol string, side OrderSide, orderType OrderType, price, amount decimal.Decimal) *Order {
	return &Order{
		ID:        id,
		Symbol:    symbol,
		Side:      side,
		Type:      orderType,
		Price:     price,
		Amount:    amount,
		Filled:    decimal.Zero,
		Remaining: amount,
		Timestamp: time.Now().UnixNano(),
		Status:    StatusNew,
	}
}

// Trade represents a successful match between two orders
type Trade struct {
	ID             string          `json:"id"`
	Symbol         string          `json:"symbol"`
	MakerByID      string          `json:"maker_order_id"`
	TakerByID      string          `json:"taker_order_id"`
	MakerAccountID string          `json:"maker_account_id"`
	TakerAccountID string          `json:"taker_account_id"`
	Price          decimal.Decimal `json:"price"`
	Amount         decimal.Decimal `json:"amount"`
	Notional       decimal.Decimal `json:"notional"`
	Fee            decimal.Decimal `json:"fee"`
	Side           string          `json:"side"` // Taker's side
	Timestamp      int64           `json:"timestamp"`
	IsBuyerMake    bool            `json:"is_buyer_maker"`
}
