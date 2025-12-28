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
	OrderStatusOpen      = "open"
	OrderStatusFilled    = "filled"
	OrderStatusCancelled = "cancelled"
)

// Subscription Tier constants
const (
	TierFree       = "Free"
	TierPro        = "Pro"
	TierEnterprise = "Enterprise"
)

// Subscription Status constants
const (
	SubscriptionStatusActive   = "active"
	SubscriptionStatusCanceled = "canceled"
	SubscriptionStatusExpired  = "expired"
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
	SubjectStrategySignal = "strategy.signal"
)

// Default settings
const (
	DefaultHistoryLimit = 100
)
