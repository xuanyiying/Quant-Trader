package biz

import (
	"context"
	"errors"

	"quant-trader/internal/model"
	"quant-trader/internal/repo"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var ErrAlertNotFound = errors.New("alert not found")

type Alert struct {
	repo   *repo.Alert
	logger *zap.Logger
}

func NewAlert(repo *repo.Alert, logger *zap.Logger) *Alert {
	return &Alert{repo: repo, logger: logger}
}

func (b *Alert) ListByUserID(ctx context.Context, userID int64) ([]model.Alert, error) {
	return b.repo.GetByUserID(ctx, userID)
}

func (b *Alert) Create(ctx context.Context, userID int64, symbol, condition string, price decimal.Decimal) (*model.Alert, error) {
	alert := &model.Alert{
		UserID:      userID,
		Symbol:      symbol,
		Condition:   condition,
		Price:       price,
		IsTriggered: false,
	}

	if err := b.repo.Create(ctx, alert); err != nil {
		return nil, err
	}

	return alert, nil
}

func (b *Alert) Delete(ctx context.Context, userID, alertID int64) error {
	alert, err := b.repo.GetByID(ctx, alertID)
	if err != nil {
		return ErrAlertNotFound
	}

	if alert.UserID != userID {
		return ErrAlertNotFound
	}

	return b.repo.Delete(ctx, alertID)
}
