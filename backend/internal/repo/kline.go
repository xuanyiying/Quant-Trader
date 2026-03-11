package repo

import (
	"context"
	"time"

	"quant-trader/internal/model"

	"gorm.io/gorm"
)

type Kline struct {
	db *gorm.DB
}

func NewKline(db *gorm.DB) *Kline {
	return &Kline{db: db}
}

func (r *Kline) GetBySymbol(ctx context.Context, symbol, period string, start, end time.Time, limit int) ([]model.KLine, error) {
	var klines []model.KLine
	query := r.db.WithContext(ctx).
		Where("symbol = ? AND period = ?", symbol, period).
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&klines).Error
	return klines, err
}

func (r *Kline) GetLatest(ctx context.Context, symbol, period string, limit int) ([]model.KLine, error) {
	var klines []model.KLine
	err := r.db.WithContext(ctx).
		Where("symbol = ? AND period = ?", symbol, period).
		Order("timestamp DESC").
		Limit(limit).
		Find(&klines).Error
	return klines, err
}

func (r *Kline) Create(ctx context.Context, kline *model.KLine) error {
	return r.db.WithContext(ctx).Create(kline).Error
}

func (r *Kline) BatchCreate(ctx context.Context, klines []model.KLine) error {
	return r.db.WithContext(ctx).CreateInBatches(klines, 100).Error
}
