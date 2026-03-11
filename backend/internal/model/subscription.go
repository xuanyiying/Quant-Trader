package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SubscriptionTier struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	Name          string          `gorm:"uniqueIndex;not null" json:"name"`
	MaxSymbols    int             `json:"max_symbols"`
	RealtimeEnabled bool          `gorm:"default:false" json:"realtime_enabled"`
	PriceMonthly  decimal.Decimal `gorm:"type:decimal(10,2)" json:"price_monthly"`
	PriceYearly   decimal.Decimal `gorm:"type:decimal(10,2)" json:"price_yearly"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (SubscriptionTier) TableName() string {
	return "subscription_tiers"
}

type UserSubscription struct {
	ID                   int64          `gorm:"primaryKey" json:"id"`
	UserID               int64          `gorm:"uniqueIndex;not null" json:"user_id"`
	TierID               int64          `gorm:"index;not null" json:"tier_id"`
	Status               string         `gorm:"not null" json:"status"`
	StripeCustomerID     string         `json:"stripe_customer_id"`
	StripeSubscriptionID string         `json:"stripe_subscription_id"`
	CurrentPeriodStart   time.Time      `json:"current_period_start"`
	CurrentPeriodEnd     time.Time      `json:"current_period_end"`
	CancelAtPeriodEnd    bool           `gorm:"default:false" json:"cancel_at_period_end"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

type SubscriptionPrice struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	TierID        int64           `gorm:"index;not null" json:"tier_id"`
	StripePriceID string          `gorm:"uniqueIndex" json:"stripe_price_id"`
	PriceAmount   decimal.Decimal `gorm:"type:decimal(10,2)" json:"price_amount"`
	BillingCycle  string          `json:"billing_cycle"`
	IsActive      bool            `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (SubscriptionPrice) TableName() string {
	return "subscription_prices"
}
