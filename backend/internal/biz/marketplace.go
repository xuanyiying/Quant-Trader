package biz

import (
	"context"

	"quant-trader/internal/repo"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Marketplace struct {
	repo   *repo.Marketplace
	logger *zap.Logger
}

func NewMarketplace(repo *repo.Marketplace, logger *zap.Logger) *Marketplace {
	return &Marketplace{repo: repo, logger: logger}
}

func (b *Marketplace) ListPublic(ctx context.Context) ([]repo.MarketItem, error) {
	return b.repo.ListPublic(ctx)
}

func (b *Marketplace) Purchase(ctx context.Context, userID, strategyID int64) error {
	price, err := b.repo.GetPrice(ctx, strategyID)
	if err != nil {
		return err
	}

	if price.GreaterThan(decimal.Zero) {
		b.logger.Info("purchasing paid strategy",
			zap.Int64("user_id", userID),
			zap.Int64("strategy_id", strategyID),
			zap.String("price", price.String()))
	}

	return b.repo.Purchase(ctx, userID, strategyID, price)
}
