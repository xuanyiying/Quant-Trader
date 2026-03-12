package model

// Side constants
const (
	SideBuy  = "buy"
	SideSell = "sell"
)

// Order Type constants
const (
	OrderTypeMarket = "market"
	OrderTypeLimit  = "limit"
)

// Order Status constants
const (
	OrderStatusNew               = "new"
	OrderStatusOpen              = "open"
	OrderStatusPartiallyFilled   = "partially_filled"
	OrderStatusFilled            = "filled"
	OrderStatusCancelled         = "cancelled"
	OrderStatusRejected          = "rejected"
	OrderStatusExpired           = "expired"
)

// Event Type constants
const (
	EventTypeSubmit      = "submit"
	EventTypeCancel      = "cancel"
	EventTypeTrigger     = "trigger"
	EventTypeTrade       = "trade"
	EventTypeOrderUpdate = "order_update"
	EventTypeReject      = "reject"
	EventTypeAlert       = "alert"
)

// Subscription Tier constants
const (
	TierFree       = "Free"
	TierPro        = "Pro"
	TierEnterprise = "Enterprise"
)

// Subscription Status constants
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusCanceled = "canceled"
	SubscriptionStatusExpired   = "expired"
	SubscriptionStatusPastDue   = "past_due"
	SubscriptionStatusTrialing  = "trialing"
	SubscriptionStatusUnpaid   = "unpaid"
)

// KLine Period constants
const (
	Period1m  = "1m"
	Period5m  = "5m"
	Period15m = "15m"
	Period1h  = "1h"
	Period4h  = "4h"
	Period1d  = "1d"
)

// Common Exchange names
const (
	ExchangeBinance = "binance"
	ExchangeBybit   = "bybit"
	ExchangeKraken  = "kraken"
)

// NATS Subject prefixes
const (
	SubjectMarketRaw      = "market.raw"
	SubjectMarketKline    = "market.kline"
	SubjectMatchingEvent  = "matching.event"
	SubjectPaperEvent     = "paper.event"
	SubjectStrategySignal = "strategy.signal"
)

// Default settings
const (
	DefaultHistoryLimit = 500
	MaxHistoryLimit     = 1000 // 系统最大限制，防止性能问题
)
