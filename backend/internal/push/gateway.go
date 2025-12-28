package push

import (
	"encoding/json"
	"net/http"
	"quant-trader/internal/infrastructure"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type PushGateway struct {
	logger        *zap.Logger
	js            nats.JetStreamContext
	clients       map[*Client]bool
	subscriptions map[string]map[*Client]bool
	natsSubs      map[string]*nats.Subscription
	mu            sync.RWMutex
}

func NewPushGateway(js nats.JetStreamContext, logger *zap.Logger) *PushGateway {
	return &PushGateway{
		logger:        logger,
		js:            js,
		clients:       make(map[*Client]bool),
		subscriptions: make(map[string]map[*Client]bool),
		natsSubs:      make(map[string]*nats.Subscription),
	}
}

func (g *PushGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		g.logger.Error("failed to upgrade websocket", zap.Error(err))
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	g.mu.Lock()
	g.clients[client] = true
	g.mu.Unlock()
	infrastructure.WSConnections.Inc()

	go g.writePump(client)
	g.readPump(client)
}

func (g *PushGateway) readPump(c *Client) {
	defer func() {
		g.mu.Lock()
		delete(g.clients, c)
		for topic, clients := range g.subscriptions {
			delete(clients, c)
			if len(clients) == 0 {
				if sub, ok := g.natsSubs[topic]; ok {
					sub.Unsubscribe()
					delete(g.natsSubs, topic)
					g.logger.Info("unsubscribed from NATS as no clients left", zap.String("topic", topic))
				}
				delete(g.subscriptions, topic)
			}
		}
		g.mu.Unlock()
		infrastructure.WSConnections.Dec()
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var req struct {
			Action string `json:"action"` // "subscribe", "unsubscribe"
			Topic  string `json:"topic"`
		}
		if err := json.Unmarshal(message, &req); err != nil {
			continue
		}

		g.mu.Lock()
		switch req.Action {
		case "subscribe":
			if g.subscriptions[req.Topic] == nil {
				g.subscriptions[req.Topic] = make(map[*Client]bool)
				if err := g.subscribeToNATS(req.Topic); err != nil {
					g.logger.Error("failed to subscribe to NATS", zap.String("topic", req.Topic), zap.Error(err))
				}
			}
			g.subscriptions[req.Topic][c] = true
			g.logger.Info("client subscribed to topic", zap.String("topic", req.Topic))
		case "unsubscribe":
			if clients, ok := g.subscriptions[req.Topic]; ok {
				delete(clients, c)
				if len(clients) == 0 {
					if sub, ok := g.natsSubs[req.Topic]; ok {
						sub.Unsubscribe()
						delete(g.natsSubs, req.Topic)
						g.logger.Info("unsubscribed from NATS as no clients left", zap.String("topic", req.Topic))
					}
					delete(g.subscriptions, req.Topic)
				}
			}
		}
		g.mu.Unlock()
	}
}

func (g *PushGateway) writePump(c *Client) {
	defer c.conn.Close()
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (g *PushGateway) subscribeToNATS(topic string) error {
	// topic can be "market.raw.*.*" or "market.kline.1m.*"
	sub, err := g.js.Subscribe(topic, func(msg *nats.Msg) {
		g.mu.RLock()
		clients := g.subscriptions[topic]
		if clients == nil {
			g.mu.RUnlock()
			return
		}

		for c := range clients {
			select {
			case c.send <- msg.Data:
			default:
				// Do not block, just drop if channel is full
			}
		}
		g.mu.RUnlock()
		msg.Ack()
	}, nats.ManualAck())

	if err != nil {
		return err
	}

	g.natsSubs[topic] = sub
	g.logger.Info("subscribed to NATS topic", zap.String("topic", topic))
	return nil
}
