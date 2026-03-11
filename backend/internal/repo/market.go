package repo

import (
	"context"

	"quant-trader/internal/model"

	"gorm.io/gorm"
)

type Market struct {
	db *gorm.DB
}

func NewMarket(db *gorm.DB) *Market {
	return &Market{db: db}
}

func (r *Market) GetSymbols(ctx context.Context) ([]string, error) {
	var symbols []string
	err := r.db.WithContext(ctx).
		Model(&model.KLine{}).
		Distinct("symbol").
		Pluck("symbol", &symbols).Error
	return symbols, err
}

func (r *Market) GetExchanges(ctx context.Context) ([]string, error) {
	var exchanges []string
	err := r.db.WithContext(ctx).
		Model(&model.KLine{}).
		Distinct("exchange").
		Pluck("exchange", &exchanges).Error
	return exchanges, err
}
