package biz

import (
	"context"
	"time"
	"quant-trader/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Backfill 负责历史数据补全的业务逻辑
type Backfill struct {
	backfiller *storage.HistoryBackfiller
	logger     *zap.Logger
}

// NewBackfill 创建 Backfill 业务逻辑实例
func NewBackfill(db *pgxpool.Pool, logger *zap.Logger) *Backfill {
	return &Backfill{
		backfiller: storage.NewHistoryBackfiller(db, logger),
		logger:     logger,
	}
}

// BackfillBinance 异步从 Binance 拉取历史 K 线数据
// 该方法会在后台 goroutine 中执行补全任务
func (b *Backfill) BackfillBinance(ctx context.Context, symbol string, startTime, endTime time.Time) error {
	// 启动异步任务
	go func() {
		backfillCtx := context.Background()
		if err := b.backfiller.BackfillBinance(backfillCtx, symbol, startTime, endTime); err != nil {
			b.logger.Error("backfill failed",
				zap.Error(err),
				zap.String("symbol", symbol),
				zap.String("exchange", "binance"),
			)
		}
	}()

	return nil
}

// BackfillExchange 支持多交易所的历史数据补全
func (b *Backfill) BackfillExchange(ctx context.Context, exchange, symbol string, startTime, endTime time.Time) error {
	switch exchange {
	case "binance":
		return b.BackfillBinance(ctx, symbol, startTime, endTime)
	default:
		b.logger.Warn("unsupported exchange for backfill",
			zap.String("exchange", exchange),
			zap.String("symbol", symbol),
		)
		return nil
	}
}
