package repo

import (
	"context"

	"quant-trader/internal/model"

	"gorm.io/gorm"
)

type Portfolio struct {
	db *gorm.DB
}

func NewPortfolio(db *gorm.DB) *Portfolio {
	return &Portfolio{db: db}
}

func (r *Portfolio) Create(ctx context.Context, p *model.Portfolio) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Portfolio) GetByID(ctx context.Context, id int64) (*model.Portfolio, error) {
	var portfolio model.Portfolio
	err := r.db.WithContext(ctx).First(&portfolio, id).Error
	return &portfolio, err
}

func (r *Portfolio) GetByUserID(ctx context.Context, userID int64) ([]model.Portfolio, error) {
	var portfolios []model.Portfolio
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&portfolios).Error
	return portfolios, err
}

func (r *Portfolio) Update(ctx context.Context, p *model.Portfolio) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *Portfolio) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Portfolio{}, id).Error
}

type PortfolioAsset struct {
	db *gorm.DB
}

func NewPortfolioAsset(db *gorm.DB) *PortfolioAsset {
	return &PortfolioAsset{db: db}
}

func (r *PortfolioAsset) GetByPortfolioID(ctx context.Context, portfolioID int64) ([]model.PortfolioAsset, error) {
	var assets []model.PortfolioAsset
	err := r.db.WithContext(ctx).Where("portfolio_id = ?", portfolioID).Find(&assets).Error
	return assets, err
}

func (r *PortfolioAsset) Upsert(ctx context.Context, asset *model.PortfolioAsset) error {
	return r.db.WithContext(ctx).Where("portfolio_id = ? AND symbol = ?", asset.PortfolioID, asset.Symbol).
		Assign(*asset).
		FirstOrCreate(asset).Error
}

func (r *PortfolioAsset) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.PortfolioAsset{}, id).Error
}
