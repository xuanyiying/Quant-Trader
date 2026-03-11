package repo

import (
	"context"
	"errors"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaperAccount struct {
	db *gorm.DB
}

func NewPaperAccount(db *gorm.DB) *PaperAccount {
	return &PaperAccount{db: db}
}

func (r *PaperAccount) GetByUserID(ctx context.Context, userID int64) (*model.PaperAccount, error) {
	var account model.PaperAccount
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *PaperAccount) GetOrCreate(ctx context.Context, userID int64, initialBalance decimal.Decimal) (*model.PaperAccount, error) {
	account, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return account, nil
	}

	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	account = &model.PaperAccount{
		UserID:         userID,
		Balance:        initialBalance,
		InitialBalance: initialBalance,
	}
	err = r.db.WithContext(ctx).Create(account).Error
	if err != nil {
		return nil, err
	}

	return r.GetByUserID(ctx, userID)
}

func (r *PaperAccount) UpdateBalance(ctx context.Context, userID int64, balance decimal.Decimal) error {
	return r.db.WithContext(ctx).Model(&model.PaperAccount{}).
		Where("user_id = ?", userID).
		Update("balance", balance).Error
}

func (r *PaperAccount) Reset(ctx context.Context, userID int64, balance decimal.Decimal) error {
	return r.db.WithContext(ctx).Model(&model.PaperAccount{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"balance":         balance,
			"initial_balance": balance,
		}).Error
}
