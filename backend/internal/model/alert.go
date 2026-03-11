package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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
