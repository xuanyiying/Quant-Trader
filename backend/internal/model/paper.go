package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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
