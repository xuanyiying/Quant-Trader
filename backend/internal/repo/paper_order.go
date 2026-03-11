package repo

import (
	"context"
	"errors"

	"quant-trader/internal/model"

	"gorm.io/gorm"
)

type PaperOrder struct {
	db *gorm.DB
}

func NewPaperOrder(db *gorm.DB) *PaperOrder {
	return &PaperOrder{db: db}
}

func (r *PaperOrder) Create(ctx context.Context, order *model.PaperOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *PaperOrder) GetByID(ctx context.Context, id int64) (*model.PaperOrder, error) {
	var order model.PaperOrder
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *PaperOrder) GetByUserID(ctx context.Context, userID int64) ([]model.PaperOrder, error) {
	var orders []model.PaperOrder
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *PaperOrder) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.PaperOrder{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *PaperOrder) CancelByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&model.PaperOrder{}).
		Where("user_id = ? AND status = ?", userID, "open").
		Update("status", "cancelled").Error
}
