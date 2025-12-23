package engine

import (
	"context"
	"market-ingestor/internal/model"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DataLoader struct {
	pool *pgxpool.Pool
}

func NewDataLoader(pool *pgxpool.Pool) *DataLoader {
	return &DataLoader{pool: pool}
}

func (l *DataLoader) LoadCandles(ctx context.Context, symbol string, start, end time.Time, period string) ([]model.KLine, error) {
	rows, err := l.pool.Query(ctx, `
		SELECT time, symbol, exchange, period, open, high, low, close, volume 
		FROM market_klines 
		WHERE symbol = $1 AND period = $2 AND time >= $3 AND time <= $4 
		ORDER BY time ASC`,
		symbol, period, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []model.KLine
	for rows.Next() {
		var k model.KLine
		if err := rows.Scan(&k.Timestamp, &k.Symbol, &k.Exchange, &k.Period, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume); err != nil {
			return nil, err
		}
		candles = append(candles, k)
	}
	return candles, nil
}
