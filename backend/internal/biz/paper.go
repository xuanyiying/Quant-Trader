package biz

import (
	"context"
	"errors"

	"quant-trader/internal/model"
	"quant-trader/internal/repo"
	"quant-trader/internal/risk"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrInvalidOrderSide = errors.New("invalid order side")
	ErrInvalidOrderType = errors.New("invalid order type")
)

type PaperTrading struct {
	paper    *repo.PaperAccount
	position *repo.PaperPosition
	order    *repo.PaperOrder
	risk     *risk.RiskManager
	logger   *zap.Logger
}

func NewPaperTrading(
	paper *repo.PaperAccount,
	position *repo.PaperPosition,
	order *repo.PaperOrder,
	risk *risk.RiskManager,
	logger *zap.Logger,
) *PaperTrading {
	return &PaperTrading{
		paper:    paper,
		position: position,
		order:    order,
		risk:     risk,
		logger:   logger,
	}
}

func (b *PaperTrading) GetOrCreateAccount(ctx context.Context, userID int64) (*model.PaperAccount, error) {
	return b.paper.GetOrCreate(ctx, userID, decimal.NewFromFloat(100000))
}

func (b *PaperTrading) ResetAccount(ctx context.Context, userID int64) error {
	initialBalance := decimal.NewFromFloat(100000)

	if err := b.paper.Reset(ctx, userID, initialBalance); err != nil {
		return err
	}

	if err := b.position.ClearAll(ctx, userID); err != nil {
		b.logger.Error("failed to clear positions after reset", zap.Error(err))
		return err
	}

	return nil
}

type CreateOrderInput struct {
	Symbol string
	Side   string
	Type   string
	Price  decimal.Decimal
	Qty    decimal.Decimal
}

func (b *PaperTrading) CreateOrder(ctx context.Context, userID int64, input CreateOrderInput) (*model.PaperOrder, error) {
	if input.Side != "buy" && input.Side != "sell" {
		return nil, ErrInvalidOrderSide
	}
	if input.Type != "market" && input.Type != "limit" {
		return nil, ErrInvalidOrderType
	}

	if err := b.risk.PreTradeCheck(ctx, userID, input.Symbol, input.Side, input.Qty, input.Price); err != nil {
		return nil, err
	}

	order := &model.PaperOrder{
		UserID: userID,
		Symbol: input.Symbol,
		Side:   input.Side,
		Type:   input.Type,
		Price:  input.Price,
		Qty:    input.Qty,
		Status: "open",
	}

	if err := b.order.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (b *PaperTrading) GetPositions(ctx context.Context, userID int64) ([]model.PaperPosition, error) {
	return b.position.GetByUserID(ctx, userID)
}
