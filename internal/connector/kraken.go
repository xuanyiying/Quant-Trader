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

type KrakenConnector struct {
	logger *zap.Logger
	symbol string // e.g. XBT/USD
}

func NewKrakenConnector(logger *zap.Logger, symbol string) *KrakenConnector {
	return &KrakenConnector{
		logger: logger,
		symbol: symbol,
	}
}

func (k *KrakenConnector) Run(ctx context.Context, tradeChan chan<- model.Trade) {
	url := "wss://ws.kraken.com"
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		k.logger.Info("connecting to Kraken websocket", zap.String("url", url))
		dialer := websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
		}
		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			k.logger.Error("failed to connect to Kraken", zap.Error(err))
			time.Sleep(backoff)
			backoff = k.increaseBackoff(backoff)
			continue
		}

		backoff = time.Second
		k.logger.Info("connected to Kraken websocket")
		infrastructure.WSConnections.Inc()

		// Subscribe
		subMsg := map[string]interface{}{
			"event": "subscribe",
			"pair": []string{
				k.symbol,
			},
			"subscription": map[string]string{
				"name": "trade",
			},
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			k.logger.Error("failed to subscribe to Kraken trades", zap.Error(err))
			conn.Close()
			continue
		}

		if err := k.handleConnection(ctx, conn, tradeChan); err != nil {
			k.logger.Error("Kraken connection closed with error", zap.Error(err))
		}
		infrastructure.WSConnections.Dec()
		conn.Close()
	}
}

func (k *KrakenConnector) handleConnection(ctx context.Context, conn *websocket.Conn, tradeChan chan<- model.Trade) error {
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Heartbeat (Kraken doesn't strictly need it if there's activity, but good practice)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteJSON(map[string]string{"event": "ping"}); err != nil {
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

			// Kraken returns an array for data, and an object for events
			if message[0] == '{' {
				// Event message (e.g. subscriptionStatus, pong, heartbeat)
				continue
			}

			var raw []interface{}
			if err := json.Unmarshal(message, &raw); err != nil {
				continue
			}

			// Expected format: [channelID, [[price, volume, time, side, orderType, misc], ...], "trade", "pair"]
			if len(raw) < 4 {
				continue
			}

			tradesData, ok := raw[1].([]interface{})
			if !ok {
				continue
			}

			pair, _ := raw[3].(string)

			for _, t := range tradesData {
				tradeArr, ok := t.([]interface{})
				if !ok || len(tradeArr) < 4 {
					continue
				}

				trade := k.convertToModel(tradeArr, pair)
				select {
				case tradeChan <- trade:
				default:
					k.logger.Warn("trade channel full, dropping Kraken trade", zap.String("trade_id", trade.ID))
				}
			}

			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		}
	}
}

func (k *KrakenConnector) convertToModel(data []interface{}, pair string) model.Trade {
	priceStr, _ := data[0].(string)
	volumeStr, _ := data[1].(string)
	timeStr, _ := data[2].(string)
	sideCode, _ := data[3].(string)

	price, _ := decimal.NewFromString(priceStr)
	volume, _ := decimal.NewFromString(volumeStr)

	// Kraken time is "1534614057.321597" (seconds.nanoseconds)
	var ts time.Time
	var sec int64
	var nsec int64
	fmt.Sscanf(timeStr, "%d.%d", &sec, &nsec)
	ts = time.Unix(sec, nsec)

	side := "buy"
	if sideCode == "s" {
		side = "sell"
	}

	return model.Trade{
		ID:        fmt.Sprintf("%d", ts.UnixNano()), // Kraken doesn't provide a unique trade ID in v1 WS
		Symbol:    pair,
		Exchange:  "kraken",
		Price:     price,
		Amount:    volume,
		Side:      side,
		Timestamp: ts,
	}
}

func (k *KrakenConnector) increaseBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > time.Minute {
		return time.Minute
	}
	return next
}
