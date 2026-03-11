package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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
