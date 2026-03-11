package biz

import (
	"context"

	"quant-trader/internal/analytics"
	"quant-trader/internal/model"
	"quant-trader/internal/repo"

	"go.uber.org/zap"
)

type Portfolio struct {
	portfolio  *repo.Portfolio
	analytics  *analytics.AnalyticsService
	logger     *zap.Logger
}

func NewPortfolio(
	portfolio *repo.Portfolio,
	analytics *analytics.AnalyticsService,
	logger *zap.Logger,
) *Portfolio {
	return &Portfolio{
		portfolio: portfolio,
		analytics: analytics,
		logger:    logger,
	}
}

func (b *Portfolio) Create(ctx context.Context, userID int64, name string) (*model.Portfolio, error) {
	p := &model.Portfolio{
		UserID: userID,
		Name:   name,
	}

	if err := b.portfolio.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (b *Portfolio) GetByUserID(ctx context.Context, userID int64) ([]model.Portfolio, error) {
	return b.portfolio.GetByUserID(ctx, userID)
}

func (b *Portfolio) GetReport(ctx context.Context, userID int64) (interface{}, error) {
	return b.analytics.GetPortfolioReport(ctx, userID)
}
