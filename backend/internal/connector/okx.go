package connector

import (
	"context"
	"encoding/json"
	"quant-trader/internal/infrastructure"
	"quant-trader/internal/model"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type OKXConnector struct {
	logger *zap.Logger
	symbol string // e.g. BTC-USDT
}

func NewOKXConnector(logger *zap.Logger, symbol string) *OKXConnector {
	return &OKXConnector{
		logger: logger,
		symbol: symbol,
	}
}

// OKXTradeEvent represents the raw trade event from OKX WS v5
type OKXTradeEvent struct {
	Arg  OKXArg         `json:"arg"`
	Data []OKXTradeData `json:"data"`
}

type OKXArg struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type OKXTradeData struct {
	InstId  string `json:"instId"`
	TradeId string `json:"tradeId"`
	Px      string `json:"px"`
	Sz      string `json:"sz"`
	Side    string `json:"side"`
	Ts      string `json:"ts"`
}

func (o *OKXConnector) Run(ctx context.Context, tradeChan chan<- model.Trade) {
	url := "wss://ws.okx.com:8443/ws/v5/public"
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		o.logger.Info("connecting to OKX websocket", zap.String("url", url))
		dialer := websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
		}
		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			o.logger.Error("failed to connect to OKX", zap.Error(err))
			time.Sleep(backoff)
			backoff = o.increaseBackoff(backoff)
			continue
		}

		backoff = time.Second // Reset backoff on successful connection
		o.logger.Info("connected to OKX websocket")
		infrastructure.WSConnections.Inc()

		// Subscribe to trades
		subMsg := map[string]interface{}{
			"op": "subscribe",
			"args": []map[string]string{
				{
					"channel": "trades",
					"instId":  o.symbol,
				},
			},
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			o.logger.Error("failed to subscribe to OKX trades", zap.Error(err))
			conn.Close()
			continue
		}

		if err := o.handleConnection(ctx, conn, tradeChan); err != nil {
			o.logger.Error("OKX connection closed with error", zap.Error(err))
		}
		infrastructure.WSConnections.Dec()
		conn.Close()
	}
}

func (o *OKXConnector) handleConnection(ctx context.Context, conn *websocket.Conn, tradeChan chan<- model.Trade) error {
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// OKX uses "ping" string for heartbeats
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
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

			// Handle pong response
			if string(message) == "pong" {
				conn.SetReadDeadline(time.Now().Add(60 * time.Second))
				continue
			}

			var event OKXTradeEvent
			if err := json.Unmarshal(message, &event); err != nil {
				o.logger.Error("failed to unmarshal OKX trade event", zap.Error(err))
				continue
			}

			// Skip subscription success messages or other non-data messages
			if len(event.Data) == 0 {
				continue
			}

			for _, data := range event.Data {
				trade := o.convertToModel(data)
				select {
				case tradeChan <- trade:
				default:
					o.logger.Warn("trade channel full, dropping OKX trade", zap.String("trade_id", trade.ID))
				}
			}
		}
	}
}

func (o *OKXConnector) convertToModel(data OKXTradeData) model.Trade {
	price, _ := decimal.NewFromString(data.Px)
	amount, _ := decimal.NewFromString(data.Sz)
	ts, _ := decimal.NewFromString(data.Ts)

	return model.Trade{
		ID:        data.TradeId,
		Symbol:    data.InstId,
		Exchange:  "okx",
		Price:     price,
		Amount:    amount,
		Side:      data.Side,
		Timestamp: time.Unix(0, ts.IntPart()*int64(time.Millisecond)),
	}
}

func (o *OKXConnector) increaseBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > time.Minute {
		return time.Minute
	}
	return next
}
