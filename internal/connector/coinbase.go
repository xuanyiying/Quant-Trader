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

type CoinbaseConnector struct {
	logger *zap.Logger
	symbol string // e.g. BTC-USD
}

func NewCoinbaseConnector(logger *zap.Logger, symbol string) *CoinbaseConnector {
	return &CoinbaseConnector{
		logger: logger,
		symbol: symbol,
	}
}

type CoinbaseMatchEvent struct {
	Type      string `json:"type"`
	TradeID   int64  `json:"trade_id"`
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"`
	Time      string `json:"time"` // RFC3339
}

func (c *CoinbaseConnector) Run(ctx context.Context, tradeChan chan<- model.Trade) {
	url := "wss://ws-feed.exchange.coinbase.com"
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		c.logger.Info("connecting to Coinbase websocket", zap.String("url", url))
		dialer := websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
		}
		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			c.logger.Error("failed to connect to Coinbase", zap.Error(err))
			time.Sleep(backoff)
			backoff = c.increaseBackoff(backoff)
			continue
		}

		backoff = time.Second
		c.logger.Info("connected to Coinbase websocket")
		infrastructure.WSConnections.Inc()

		// Subscribe
		subMsg := map[string]interface{}{
			"type": "subscribe",
			"channels": []map[string]interface{}{
				{
					"name": "matches",
					"product_ids": []string{
						c.symbol,
					},
				},
			},
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			c.logger.Error("failed to subscribe to Coinbase trades", zap.Error(err))
			conn.Close()
			continue
		}

		if err := c.handleConnection(ctx, conn, tradeChan); err != nil {
			c.logger.Error("Coinbase connection closed with error", zap.Error(err))
		}
		infrastructure.WSConnections.Dec()
		conn.Close()
	}
}

func (c *CoinbaseConnector) handleConnection(ctx context.Context, conn *websocket.Conn, tradeChan chan<- model.Trade) error {
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Coinbase doesn't require explicit ping, but we can send one if needed.
	// Actually, they recommend sending a heartbeat or just relying on the feed.

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			var event CoinbaseMatchEvent
			if err := json.Unmarshal(message, &event); err != nil {
				continue
			}

			if event.Type != "match" && event.Type != "last_match" {
				continue
			}

			trade := c.convertToModel(event)
			select {
			case tradeChan <- trade:
			default:
				c.logger.Warn("trade channel full, dropping Coinbase trade", zap.String("trade_id", trade.ID))
			}

			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		}
	}
}

func (c *CoinbaseConnector) convertToModel(event CoinbaseMatchEvent) model.Trade {
	price, _ := decimal.NewFromString(event.Price)
	amount, _ := decimal.NewFromString(event.Size)
	t, _ := time.Parse(time.RFC3339, event.Time)

	return model.Trade{
		ID:        fmt.Sprintf("%d", event.TradeID),
		Symbol:    event.ProductID,
		Exchange:  "coinbase",
		Price:     price,
		Amount:    amount,
		Side:      event.Side,
		Timestamp: t,
	}
}

func (c *CoinbaseConnector) increaseBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > time.Minute {
		return time.Minute
	}
	return next
}
