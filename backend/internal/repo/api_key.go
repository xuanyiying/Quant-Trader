package repo

import (
	"context"
	"errors"

	"quant-trader/internal/model"

	"gorm.io/gorm"
)

type APIKey struct {
	db *gorm.DB
}

func NewAPIKey(db *gorm.DB) *APIKey {
	return &APIKey{db: db}
}

func (r *APIKey) Create(ctx context.Context, key *model.APIKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}

func (r *APIKey) GetByID(ctx context.Context, id int64) (*model.APIKey, error) {
	var key model.APIKey
	err := r.db.WithContext(ctx).First(&key, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &key, nil
}

func (r *APIKey) GetByUserID(ctx context.Context, userID int64) ([]model.APIKey, error) {
	var keys []model.APIKey
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&keys).Error
	return keys, err
}

func (r *APIKey) GetByKeyHash(ctx context.Context, keyHash string) (*model.APIKey, error) {
	var key model.APIKey
	err := r.db.WithContext(ctx).Where("key_hash = ?", keyHash).First(&key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &key, nil
}

func (r *APIKey) Update(ctx context.Context, key *model.APIKey) error {
	return r.db.WithContext(ctx).Save(key).Error
}

func (r *APIKey) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.APIKey{}, id).Error
}

func (r *APIKey) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.APIKey{}).Error
}
