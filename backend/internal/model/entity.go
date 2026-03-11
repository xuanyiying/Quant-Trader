package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Nickname     string         `json:"nickname"`
	TierID       int64          `gorm:"default:1" json:"tier_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type UserExchangeKey struct {
	ID         int64          `gorm:"primaryKey" json:"id"`
	UserID     int64          `gorm:"index;not null" json:"user_id"`
	Exchange   string         `gorm:"not null" json:"exchange"`
	APIKey     string         `gorm:"not null" json:"api_key"`
	SecretKey  string         `gorm:"not null" json:"-"`
	Passphrase string         `json:"passphrase,omitempty"`
	Label      string         `json:"label"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserExchangeKey) TableName() string {
	return "user_exchange_keys"
}

type APIKey struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	UserID      int64          `gorm:"index;not null" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	KeyHash     string         `gorm:"not null" json:"-"`
	KeyPrefix   string         `gorm:"not null" json:"key_prefix"`
	Permissions string         `json:"permissions"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (APIKey) TableName() string {
	return "api_keys"
}

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

type Alert struct {
	ID          int64           `gorm:"primaryKey" json:"id"`
	UserID      int64           `gorm:"index;not null" json:"user_id"`
	Symbol      string          `gorm:"not null" json:"symbol"`
	Condition   string          `gorm:"not null" json:"condition"`
	Price       decimal.Decimal `gorm:"type:decimal(20,8)" json:"price"`
	IsTriggered bool            `gorm:"default:false" json:"is_triggered"`
	TriggeredAt *time.Time      `json:"triggered_at"`
	Message     string          `json:"message"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Alert) TableName() string {
	return "alerts"
}

type PaperAccount struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	UserID         int64           `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance        decimal.Decimal `gorm:"type:decimal(20,8);default:0" json:"balance"`
	InitialBalance decimal.Decimal `gorm:"type:decimal(20,8)" json:"initial_balance"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (PaperAccount) TableName() string {
	return "paper_accounts"
}

type PaperPosition struct {
	ID        int64           `gorm:"primaryKey" json:"id"`
	UserID    int64           `gorm:"index;not null" json:"user_id"`
	Symbol    string          `gorm:"index;not null" json:"symbol"`
	Side      string          `gorm:"not null" json:"side"`
	Qty       decimal.Decimal `gorm:"type:decimal(20,8)" json:"qty"`
	AvgPrice  decimal.Decimal `gorm:"type:decimal(20,8)" json:"avg_price"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (PaperPosition) TableName() string {
	return "paper_positions"
}

type PaperOrder struct {
	ID        int64           `gorm:"primaryKey" json:"id"`
	UserID    int64           `gorm:"index;not null" json:"user_id"`
	Symbol    string          `gorm:"index;not null" json:"symbol"`
	Side      string          `gorm:"not null" json:"side"`
	Type      string          `gorm:"not null" json:"type"`
	Price     decimal.Decimal `gorm:"type:decimal(20,8)" json:"price"`
	Qty       decimal.Decimal `gorm:"type:decimal(20,8)" json:"qty"`
	FilledQty decimal.Decimal `gorm:"type:decimal(20,8);default:0" json:"filled_qty"`
	Status    string          `gorm:"default:'open'" json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (PaperOrder) TableName() string {
	return "paper_orders"
}

type Portfolio struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	UserID    int64          `gorm:"index;not null" json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Portfolio) TableName() string {
	return "portfolios"
}

type PortfolioAsset struct {
	ID          int64           `gorm:"primaryKey" json:"id"`
	PortfolioID int64           `gorm:"index;not null" json:"portfolio_id"`
	Symbol      string          `gorm:"not null" json:"symbol"`
	Side        string          `gorm:"not null" json:"side"`
	Qty         decimal.Decimal `gorm:"type:decimal(20,8)" json:"qty"`
	EntryPrice  decimal.Decimal `gorm:"type:decimal(20,8)" json:"entry_price"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (PortfolioAsset) TableName() string {
	return "portfolio_assets"
}

type StrategyMarket struct {
	ID          int64           `gorm:"primaryKey" json:"id"`
	UserID      int64           `gorm:"index;not null" json:"user_id"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2)" json:"price"`
	IsPublished bool            `gorm:"default:false" json:"is_published"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (StrategyMarket) TableName() string {
	return "strategy_market"
}

type StrategyPurchase struct {
	ID         int64           `gorm:"primaryKey" json:"id"`
	UserID     int64           `gorm:"index;not null" json:"user_id"`
	StrategyID int64           `gorm:"index;not null" json:"strategy_id"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2)" json:"price"`
	LicenseKey string          `json:"license_key"`
	CreatedAt  time.Time       `json:"created_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (StrategyPurchase) TableName() string {
	return "strategy_purchases"
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
