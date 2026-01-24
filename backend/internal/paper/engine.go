package paper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"quant-trader/internal/model"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Order struct {
	ID          int64
	UserID      int64
	Symbol      string
	Side        string
	Type        string
	Price       decimal.Decimal
	Qty         decimal.Decimal
	Status      string
	FilledPrice decimal.Decimal
}

type PaperEngine struct {
	db           *pgxpool.Pool
	js           nats.JetStreamContext
	logger       *zap.Logger
	orders       map[string][]Order // key: symbol
	latestPrices map[string]decimal.Decimal
	mu           sync.RWMutex
	fillChan     chan Order
}

func NewPaperEngine(db *pgxpool.Pool, js nats.JetStreamContext, logger *zap.Logger) *PaperEngine {
	return &PaperEngine{
		db:           db,
		js:           js,
		logger:       logger,
		orders:       make(map[string][]Order),
		latestPrices: make(map[string]decimal.Decimal),
		fillChan:     make(chan Order, 1000),
	}
}

func (e *PaperEngine) Start(ctx context.Context) error {
	// 1. Load open orders from DB
	if err := e.loadOpenOrders(ctx); err != nil {
		return err
	}

	// 2. Subscribe to price updates
	_, err := e.js.Subscribe(fmt.Sprintf("%s.%s.*", model.SubjectMarketKline, model.Period1m), func(msg *nats.Msg) {
		var candle model.KLine
		if err := json.Unmarshal(msg.Data, &candle); err != nil {
			return
		}
		e.processPriceUpdate(candle)
	})

	if err != nil {
		return err
	}

	go e.batchFlushLoop(ctx)
	e.logger.Info("paper trading engine started with batching")
	return nil
}

func (e *PaperEngine) loadOpenOrders(ctx context.Context) error {
	rows, err := e.db.Query(ctx, "SELECT id, user_id, symbol, side, type, price, qty, status FROM paper_orders WHERE status = $1", model.OrderStatusOpen)
	if err != nil {
		return err
	}
	defer rows.Close()

	e.mu.Lock()
	defer e.mu.Unlock()
	e.orders = make(map[string][]Order)

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Qty, &o.Status); err != nil {
			continue
		}
		e.orders[o.Symbol] = append(e.orders[o.Symbol], o)
	}
	return nil
}

func (e *PaperEngine) processPriceUpdate(candle model.KLine) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.latestPrices[candle.Symbol] = candle.Close
	orders, ok := e.orders[candle.Symbol]
	if !ok || len(orders) == 0 {
		return
	}

	var remaining []Order
	var toFill []Order

	for _, o := range orders {
		filled := false
		switch o.Type {
		case model.OrderTypeMarket:
			filled = true
			o.FilledPrice = candle.Close
		case model.OrderTypeLimit:
			if o.Side == model.SideBuy && candle.Low.LessThanOrEqual(o.Price) {
				filled = true
				o.FilledPrice = o.Price
			} else if o.Side == model.SideSell && candle.High.GreaterThanOrEqual(o.Price) {
				filled = true
				o.FilledPrice = o.Price
			}
		}

		if filled {
			toFill = append(toFill, o)
		} else {
			remaining = append(remaining, o)
		}
	}

	if len(toFill) > 0 {
		e.orders[candle.Symbol] = remaining
	}

	for _, o := range toFill {
		e.fillChan <- o
	}
}

func (e *PaperEngine) batchFlushLoop(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var batch []Order
	for {
		select {
		case <-ctx.Done():
			return
		case o := <-e.fillChan:
			batch = append(batch, o)
			if len(batch) >= 50 {
				e.flushBatch(ctx, batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				e.flushBatch(ctx, batch)
				batch = nil
			}
		}
	}
}

func (e *PaperEngine) flushBatch(ctx context.Context, batch []Order) {
	tx, err := e.db.Begin(ctx)
	if err != nil {
		e.logger.Error("failed to start batch transaction", zap.Error(err))
		return
	}
	defer tx.Rollback(ctx)

	for _, o := range batch {
		// 1. Update order status
		_, err = tx.Exec(ctx, fmt.Sprintf("UPDATE paper_orders SET status = '%s', filled_price = $1, updated_at = NOW() WHERE id = $2", model.OrderStatusFilled),
			o.FilledPrice, o.ID)
		if err != nil {
			continue
		}

		// 2. Update balance
		amount := o.Qty.Mul(o.FilledPrice)
		if o.Side == model.SideBuy {
			_, err = tx.Exec(ctx, "UPDATE paper_accounts SET balance = balance - $1 WHERE user_id = $2", amount, o.UserID)
		} else {
			_, err = tx.Exec(ctx, "UPDATE paper_accounts SET balance = balance + $1 WHERE user_id = $2", amount, o.UserID)
		}
		if err != nil {
			continue
		}

		// 3. Update positions
		if o.Side == model.SideBuy {
			_, err = tx.Exec(ctx, `
				INSERT INTO paper_positions (user_id, symbol, qty, avg_price) 
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (user_id, symbol) DO UPDATE SET
					avg_price = (paper_positions.qty * paper_positions.avg_price + $3 * $4) / (paper_positions.qty + $3),
					qty = paper_positions.qty + $3`,
				o.UserID, o.Symbol, o.Qty, o.FilledPrice)
		} else {
			_, err = tx.Exec(ctx, "UPDATE paper_positions SET qty = qty - $1 WHERE user_id = $2 AND symbol = $3",
				o.Qty, o.UserID, o.Symbol)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		e.logger.Error("failed to commit batch transaction", zap.Error(err))
		return
	}

	// Publish events after successful commit
	for _, o := range batch {
		e.publishEvent(model.EventTypeTrade, o)
	}

	e.logger.Info("paper order batch flushed", zap.Int("count", len(batch)))
}

func (e *PaperEngine) publishEvent(eventType string, o Order) {
	ev := model.Event{
		Type:      eventType,
		Symbol:    o.Symbol,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Order: &model.Order{
			ID:          fmt.Sprintf("%d", o.ID),
			UserID:      o.UserID,
			Symbol:      o.Symbol,
			Side:        o.Side,
			Type:        o.Type,
			Price:       o.Price,
			Qty:         o.Qty,
			FilledQty:   o.Qty, // For now assume full fill in paper
			FilledPrice: o.FilledPrice,
			Status:      o.Status,
			Timestamp:   time.Now().UnixNano(),
		},
	}

	if eventType == model.EventTypeTrade {
		ev.Execution = &model.Execution{
			ID:        fmt.Sprintf("exec-%d-%d", o.ID, ev.Timestamp),
			Symbol:    o.Symbol,
			Price:     o.FilledPrice,
			Qty:       o.Qty,
			Side:      o.Side,
			Timestamp: ev.Timestamp,
		}
		ev.Order.Status = model.OrderStatusFilled
	}

	subject := fmt.Sprintf("%s.%s", model.SubjectPaperEvent, o.Symbol)
	data, _ := json.Marshal(ev)
	_, err := e.js.Publish(subject, data)
	if err != nil {
		e.logger.Error("failed to publish paper event", zap.Error(err))
	}
}

func (e *PaperEngine) PlaceOrder(ctx context.Context, o Order) (int64, error) {
	// Simple validation: check balance if buy
	if o.Side == model.SideBuy {
		var balance decimal.Decimal
		err := e.db.QueryRow(ctx, "SELECT balance FROM paper_accounts WHERE user_id = $1", o.UserID).Scan(&balance)
		if err != nil {
			return 0, fmt.Errorf("failed to get balance: %w", err)
		}
		cost := o.Qty.Mul(o.Price)
		if o.Type == model.OrderTypeMarket {
			e.mu.RLock()
			price, ok := e.latestPrices[o.Symbol]
			e.mu.RUnlock()
			if ok {
				cost = o.Qty.Mul(price)
			} else {
				// If no price data yet, we can't accurately estimate, but 0 is too risky.
				// Maybe return error or use a very high estimate.
				return 0, errors.New("no market price available for market order estimation")
			}
		}
		if balance.LessThan(cost) {
			return 0, errors.New("insufficient balance")
		}
	}

	var id int64
	err := e.db.QueryRow(ctx,
		"INSERT INTO paper_orders (user_id, symbol, side, type, price, qty, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		o.UserID, o.Symbol, o.Side, o.Type, o.Price, o.Qty, model.OrderStatusOpen).Scan(&id)

	if err != nil {
		return 0, err
	}

	o.ID = id
	e.mu.Lock()
	e.orders[o.Symbol] = append(e.orders[o.Symbol], o)
	e.mu.Unlock()

	e.publishEvent(model.EventTypeSubmit, o)

	return id, nil
}
