package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"market-ingestor/internal/infrastructure"
	"market-ingestor/internal/model"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type BinanceConnector struct {
	logger *zap.Logger
	symbol string
}

func NewBinanceConnector(logger *zap.Logger, symbol string) *BinanceConnector {
	return &BinanceConnector{
		logger: logger,
		symbol: symbol,
	}
}

// BinanceTradeEvent represents the raw trade event from Binance WS
type BinanceTradeEvent struct {
	EventType    string `json:"e"`
	EventTime    int64  `json:"E"`
	Symbol       string `json:"s"`
	TradeID      int64  `json:"t"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	BuyerID      int64  `json:"b"`
	SellerID     int64  `json:"a"`
	TradeTime    int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
	Ignore       bool   `json:"M"`
}

func (b *BinanceConnector) Run(ctx context.Context, tradeChan chan<- model.Trade) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@trade", b.symbol)
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		b.logger.Info("connecting to binance websocket", zap.String("url", url))
		dialer := websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
		}
		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			b.logger.Error("failed to connect to binance", zap.Error(err))
			time.Sleep(backoff)
			backoff = b.increaseBackoff(backoff)
			continue
		}

		backoff = time.Second // Reset backoff on successful connection
		b.logger.Info("connected to binance websocket")
		infrastructure.WSConnections.Inc()

		if err := b.handleConnection(ctx, conn, tradeChan); err != nil {
			b.logger.Error("connection closed with error", zap.Error(err))
		}
		infrastructure.WSConnections.Dec()
		conn.Close()
	}
}

func (b *BinanceConnector) handleConnection(ctx context.Context, conn *websocket.Conn, tradeChan chan<- model.Trade) error {
	// Ping/Pong is handled automatically by gorilla/websocket default handlers if we don't override them.
	// But we can set a read deadline to detect stale connections.
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			var event BinanceTradeEvent
			if err := json.Unmarshal(message, &event); err != nil {
				b.logger.Error("failed to unmarshal binance trade event", zap.Error(err))
				continue
			}

			trade := b.convertToModel(event)
			select {
			case tradeChan <- trade:
			default:
				b.logger.Warn("trade channel full, dropping trade", zap.String("trade_id", trade.ID))
			}
		}
	}
}

func (b *BinanceConnector) convertToModel(event BinanceTradeEvent) model.Trade {
	price, _ := decimal.NewFromString(event.Price)
	amount, _ := decimal.NewFromString(event.Quantity)

	side := "buy"
	if event.IsBuyerMaker {
		side = "sell" // Maker is buyer means taker is seller
	}

	return model.Trade{
		ID:        fmt.Sprintf("%d", event.TradeID),
		Symbol:    event.Symbol,
		Exchange:  "binance",
		Price:     price,
		Amount:    amount,
		Side:      side,
		Timestamp: time.Unix(0, event.TradeTime*int64(time.Millisecond)),
	}
}

func (b *BinanceConnector) increaseBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > time.Minute {
		return time.Minute
	}
	return next
}
