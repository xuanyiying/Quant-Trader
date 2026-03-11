package biz

import (
	"context"
	"errors"

	"quant-trader/internal/model"
	"quant-trader/internal/repo"

	"go.uber.org/zap"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type Subscription struct {
	subscription *repo.Subscription
	logger       *zap.Logger
}

func NewSubscription(subscription *repo.Subscription, logger *zap.Logger) *Subscription {
	return &Subscription{
		subscription: subscription,
		logger:       logger,
	}
}

type SubscriptionInfo struct {
	TierName   string
	MaxSymbols int
	Status     string
	ExpiresAt  string
}

func (b *Subscription) GetSubscription(ctx context.Context, userID int64) (*SubscriptionInfo, error) {
	sub, err := b.subscription.GetUserSubscription(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &SubscriptionInfo{
				TierName:   model.TierFree,
				MaxSymbols: 1,
				Status:     model.SubscriptionStatusActive,
			}, nil
		}
		return nil, err
	}

	tier, err := b.subscription.GetTierByID(ctx, sub.TierID)
	if err != nil {
		return nil, err
	}

	return &SubscriptionInfo{
		TierName:   tier.Name,
		MaxSymbols: tier.MaxSymbols,
		Status:     sub.Status,
		ExpiresAt:  sub.CurrentPeriodEnd.Format("2006-01-02"),
	}, nil
}
