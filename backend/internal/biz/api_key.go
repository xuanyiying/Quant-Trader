package biz

import (
	"context"
	"fmt"
	"time"

	"quant-trader/internal/model"
	"quant-trader/internal/repo"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type APIKey struct {
	repo   *repo.APIKey
	logger *zap.Logger
}

func NewAPIKey(repo *repo.APIKey, logger *zap.Logger) *APIKey {
	return &APIKey{repo: repo, logger: logger}
}

func (b *APIKey) ListByUserID(ctx context.Context, userID int64) ([]model.APIKey, error) {
	return b.repo.GetByUserID(ctx, userID)
}

type CreateAPIKeyResult struct {
	KeyID     string
	KeySecret string
}

func (b *APIKey) Create(ctx context.Context, userID int64, name string) (*CreateAPIKeyResult, error) {
	keyID := fmt.Sprintf("qt_%d_%d", userID, time.Now().Unix())
	rawSecret := fmt.Sprintf("sec_%d_%d", userID, time.Now().UnixNano())
	hash, err := bcrypt.GenerateFromPassword([]byte(rawSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	key := &model.APIKey{
		UserID:    userID,
		Name:      name,
		KeyPrefix: keyID,
		KeyHash:   string(hash),
		IsActive:  true,
	}

	if err := b.repo.Create(ctx, key); err != nil {
		return nil, err
	}

	return &CreateAPIKeyResult{
		KeyID:     keyID,
		KeySecret: rawSecret,
	}, nil
}
