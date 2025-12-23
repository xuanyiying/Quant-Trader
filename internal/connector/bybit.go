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

type BybitConnector struct {
	logger *zap.Logger
	symbol string // e.g. BTCUSDT
}

func NewBybitConnector(logger *zap.Logger, symbol string) *BybitConnector {
	return &BybitConnector{
		logger: logger,
		symbol: symbol,
	}
}

type BybitTradeEvent struct {
	Topic string           `json:"topic"`
	Type  string           `json:"type"`
	Ts    int64            `json:"ts"`
	Data  []BybitTradeData `json:"data"`
}

type BybitTradeData struct {
	T  int64  `json:"T"`
	S  string `json:"s"`
	S2 string `json:"S"` // Side: Buy/Sell
	P  string `json:"p"`
	V  string `json:"v"`
	I  string `json:"i"` // Trade ID
	L  string `json:"L"` // Tick direction
	B  bool   `json:"B"` // Is block trade
}

func (b *BybitConnector) Run(ctx context.Context, tradeChan chan<- model.Trade) {
	url := "wss://stream.bybit.com/v5/public/spot"
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		b.logger.Info("connecting to Bybit websocket", zap.String("url", url))
		dialer := websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
		}
		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			b.logger.Error("failed to connect to Bybit", zap.Error(err))
			time.Sleep(backoff)
			backoff = b.increaseBackoff(backoff)
			continue
		}

		backoff = time.Second
		b.logger.Info("connected to Bybit websocket")
		infrastructure.WSConnections.Inc()

		// Subscribe
		subMsg := map[string]interface{}{
			"op": "subscribe",
			"args": []string{
				fmt.Sprintf("publicTrade.%s", b.symbol),
			},
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			b.logger.Error("failed to subscribe to Bybit trades", zap.Error(err))
			conn.Close()
			continue
		}

		if err := b.handleConnection(ctx, conn, tradeChan); err != nil {
			b.logger.Error("Bybit connection closed with error", zap.Error(err))
		}
		infrastructure.WSConnections.Dec()
		conn.Close()
	}
}

func (b *BybitConnector) handleConnection(ctx context.Context, conn *websocket.Conn, tradeChan chan<- model.Trade) error {
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Heartbeat
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteJSON(map[string]string{"op": "ping"}); err != nil {
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			var event BybitTradeEvent
			if err := json.Unmarshal(message, &event); err != nil {
				// Might be pong or subscription response
				continue
			}

			if event.Topic == "" || len(event.Data) == 0 {
				continue
			}

			for _, data := range event.Data {
				trade := b.convertToModel(data)
				select {
				case tradeChan <- trade:
				default:
					b.logger.Warn("trade channel full, dropping Bybit trade", zap.String("trade_id", trade.ID))
				}
			}
		}
	}
}

func (b *BybitConnector) convertToModel(data BybitTradeData) model.Trade {
	price, _ := decimal.NewFromString(data.P)
	amount, _ := decimal.NewFromString(data.V)

	side := "buy"
	if data.S2 == "Sell" {
		side = "sell"
	}

	return model.Trade{
		ID:        data.I,
		Symbol:    data.S,
		Exchange:  "bybit",
		Price:     price,
		Amount:    amount,
		Side:      side,
		Timestamp: time.Unix(0, data.T*int64(time.Millisecond)),
	}
}

func (b *BybitConnector) increaseBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > time.Minute {
		return time.Minute
	}
	return next
}
