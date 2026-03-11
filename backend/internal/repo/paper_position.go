package repo

import (
	"context"
	"errors"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaperPosition struct {
	db *gorm.DB
}

func NewPaperPosition(db *gorm.DB) *PaperPosition {
	return &PaperPosition{db: db}
}

func (r *PaperPosition) GetByUserID(ctx context.Context, userID int64) ([]model.PaperPosition, error) {
	var positions []model.PaperPosition
	err := r.db.WithContext(ctx).Where("user_id = ? AND qty > 0", userID).Find(&positions).Error
	return positions, err
}

func (r *PaperPosition) GetBySymbol(ctx context.Context, userID int64, symbol string) (*model.PaperPosition, error) {
	var pos model.PaperPosition
	err := r.db.WithContext(ctx).Where("user_id = ? AND symbol = ?", userID, symbol).First(&pos).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pos, nil
}

func (r *PaperPosition) Upsert(ctx context.Context, userID int64, symbol, side string, qty, avgPrice decimal.Decimal) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND symbol = ?", userID, symbol).
		Assign(model.PaperPosition{
			UserID:   userID,
			Symbol:   symbol,
			Side:     side,
			Qty:      qty,
			AvgPrice: avgPrice,
		}).
		FirstOrCreate(&model.PaperPosition{}).Error
}

func (r *PaperPosition) Delete(ctx context.Context, userID int64, symbol string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND symbol = ?", userID, symbol).Delete(&model.PaperPosition{}).Error
}

func (r *PaperPosition) ClearAll(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.PaperPosition{}).Error
}
