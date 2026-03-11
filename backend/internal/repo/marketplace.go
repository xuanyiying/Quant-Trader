package repo

import (
	"context"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Marketplace struct {
	db *gorm.DB
}

func NewMarketplace(db *gorm.DB) *Marketplace {
	return &Marketplace{db: db}
}

type MarketItem struct {
	ID          int64
	Name        string
	Description string
	Price       decimal.Decimal
	AuthorEmail string
	Metrics     string
}

func (r *Marketplace) ListPublic(ctx context.Context) ([]MarketItem, error) {
	var items []MarketItem
	err := r.db.WithContext(ctx).
		Table("strategy_market m").
		Select("m.id, m.price, m.description, m.performance_metrics as metrics, s.name, u.email as author_email").
		Joins("JOIN strategies s ON m.strategy_id = s.id").
		Joins("JOIN users u ON m.owner_id = u.id").
		Where("m.is_public = ?", true).
		Scan(&items).Error
	return items, err
}

func (r *Marketplace) GetPrice(ctx context.Context, id int64) (decimal.Decimal, error) {
	var item model.StrategyMarket
	err := r.db.WithContext(ctx).
		Select("price").
		Where("id = ? AND is_published = ?", id, true).
		First(&item).Error
	return item.Price, err
}

func (r *Marketplace) Purchase(ctx context.Context, userID, strategyID int64, price decimal.Decimal) error {
	purchase := &model.StrategyPurchase{
		UserID:     userID,
		StrategyID: strategyID,
		Price:      price,
	}
	return r.db.WithContext(ctx).Create(purchase).Error
}
