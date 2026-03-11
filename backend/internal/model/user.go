package model

import (
	"time"

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
