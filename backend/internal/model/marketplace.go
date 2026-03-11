package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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
