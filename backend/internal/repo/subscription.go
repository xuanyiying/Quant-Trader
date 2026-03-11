package repo

import (
	"context"

	"quant-trader/internal/model"

	"gorm.io/gorm"
)

type Subscription struct {
	db *gorm.DB
}

func NewSubscription(db *gorm.DB) *Subscription {
	return &Subscription{db: db}
}

func (r *Subscription) GetTierByID(ctx context.Context, id int64) (*model.SubscriptionTier, error) {
	var tier model.SubscriptionTier
	err := r.db.WithContext(ctx).First(&tier, id).Error
	return &tier, err
}

func (r *Subscription) GetTierByName(ctx context.Context, name string) (*model.SubscriptionTier, error) {
	var tier model.SubscriptionTier
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tier).Error
	return &tier, err
}

func (r *Subscription) GetUserSubscription(ctx context.Context, userID int64) (*model.UserSubscription, error) {
	var sub model.UserSubscription
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&sub).Error
	return &sub, err
}

func (r *Subscription) CreateUserSubscription(ctx context.Context, sub *model.UserSubscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *Subscription) UpdateUserSubscription(ctx context.Context, sub *model.UserSubscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *Subscription) GetActivePriceIDs(ctx context.Context) ([]model.SubscriptionPrice, error) {
	var prices []model.SubscriptionPrice
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&prices).Error
	return prices, err
}

func (r *Subscription) GetTierPrices(ctx context.Context, tierID int64) ([]model.SubscriptionPrice, error) {
	var prices []model.SubscriptionPrice
	err := r.db.WithContext(ctx).Where("tier_id = ? AND is_active = ?", tierID, true).Find(&prices).Error
	return prices, err
}
