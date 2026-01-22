package matching

import (
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type AbuseConfig struct {
	MaxOrdersPerSecond   int
	CancelRatioWindowSec int64
	MaxCancelRatioBps    int64
	SpoofNotional        decimal.Decimal
	SpoofWindowMs        int64
	SpoofMaxCount        int
	BlockSelfTrade       bool
	BlockDurationMs      int64
}

type AbuseRiskManager struct {
	mu          sync.Mutex
	cfg         AbuseConfig
	orders      map[string]map[int64]int
	cancels     map[string]map[int64]int
	blockedTill map[string]int64
	orderMeta   map[string]abuseOrderMeta
}

type abuseOrderMeta struct {
	AccountID string
	PlacedAt  int64
	Notional  decimal.Decimal
}

func NewAbuseRiskManager(cfg AbuseConfig) *AbuseRiskManager {
	if cfg.MaxOrdersPerSecond <= 0 {
		cfg.MaxOrdersPerSecond = 200
	}
	if cfg.CancelRatioWindowSec <= 0 {
		cfg.CancelRatioWindowSec = 60
	}
	if cfg.MaxCancelRatioBps <= 0 {
		cfg.MaxCancelRatioBps = 9000
	}
	if cfg.SpoofNotional.LessThanOrEqual(decimal.Zero) {
		cfg.SpoofNotional = decimal.RequireFromString("100000")
	}
	if cfg.SpoofWindowMs <= 0 {
		cfg.SpoofWindowMs = 500
	}
	if cfg.SpoofMaxCount <= 0 {
		cfg.SpoofMaxCount = 20
	}
	if cfg.BlockDurationMs <= 0 {
		cfg.BlockDurationMs = int64((30 * time.Second) / time.Millisecond)
	}
	return &AbuseRiskManager{
		cfg:         cfg,
		orders:      make(map[string]map[int64]int),
		cancels:     make(map[string]map[int64]int),
		blockedTill: make(map[string]int64),
		orderMeta:   make(map[string]abuseOrderMeta),
	}
}

func (r *AbuseRiskManager) ValidateOrder(state MarketState, lastPrice decimal.Decimal, snap OrderBookSnapshot, order *Order) error {
	account := order.AccountID
	if account == "" {
		account = "_"
	}
	ts := order.Timestamp
	if ts <= 0 {
		ts = time.Now().UnixNano()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if until, ok := r.blockedTill[account]; ok && until > ts {
		return fmt.Errorf("account blocked")
	}

	sec := ts / 1_000_000_000
	if r.inc(r.orders, account, sec, r.cfg.MaxOrdersPerSecond) {
		return fmt.Errorf("account rate limit exceeded")
	}

	windowStart := sec - r.cfg.CancelRatioWindowSec + 1
	orders := r.sumSince(r.orders[account], windowStart)
	cancels := r.sumSince(r.cancels[account], windowStart)
	if orders > 0 {
		ratioBps := int64(cancels) * 10000 / int64(orders)
		if ratioBps > r.cfg.MaxCancelRatioBps {
			return fmt.Errorf("cancel ratio too high")
		}
	}

	return nil
}

func (r *AbuseRiskManager) OnTrade(trade Trade) {
	if !r.cfg.BlockSelfTrade {
		return
	}
	if trade.MakerAccountID == "" || trade.TakerAccountID == "" {
		return
	}
	if trade.MakerAccountID != trade.TakerAccountID {
		return
	}
	now := time.Now().UnixNano()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blockedTill[trade.MakerAccountID] = now + r.cfg.BlockDurationMs*1_000_000
}

func (r *AbuseRiskManager) ConsumeEvent(ev Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch ev.Type {
	case EventSubmit:
		if ev.Order == nil {
			return
		}
		account := ev.Order.AccountID
		if account == "" {
			account = "_"
		}
		r.orderMeta[ev.Order.ID] = abuseOrderMeta{
			AccountID: account,
			PlacedAt:  ev.Timestamp,
			Notional:  ev.Order.Price.Mul(ev.Order.Amount),
		}
	case EventOrderUpdate:
		if ev.Order == nil {
			return
		}
		account := ev.Order.AccountID
		if account == "" {
			account = "_"
		}
		sec := ev.Timestamp / 1_000_000_000
		if ev.Order.Status == StatusCanceled || ev.Order.Status == StatusPartiallyFilledCanceled {
			r.inc(r.cancels, account, sec, 0)

			meta, ok := r.orderMeta[ev.Order.ID]
			if ok {
				dt := ev.Timestamp - meta.PlacedAt
				if dt >= 0 && dt <= r.cfg.SpoofWindowMs*1_000_000 && meta.Notional.GreaterThanOrEqual(r.cfg.SpoofNotional) {
					r.blockedTill[account] = ev.Timestamp + r.cfg.BlockDurationMs*1_000_000
				}
				delete(r.orderMeta, ev.Order.ID)
			}
		}
	case EventCancel:
		if ev.OrderID == "" {
			return
		}
		meta, ok := r.orderMeta[ev.OrderID]
		if !ok {
			return
		}
		account := meta.AccountID
		sec := ev.Timestamp / 1_000_000_000
		r.inc(r.cancels, account, sec, 0)
	case EventTrade:
		if ev.Trade == nil {
			return
		}
		if r.cfg.BlockSelfTrade && ev.Trade.MakerAccountID != "" && ev.Trade.MakerAccountID == ev.Trade.TakerAccountID {
			r.blockedTill[ev.Trade.MakerAccountID] = ev.Timestamp + r.cfg.BlockDurationMs*1_000_000
		}
	}
}

func (r *AbuseRiskManager) inc(m map[string]map[int64]int, account string, sec int64, limit int) bool {
	b, ok := m[account]
	if !ok {
		b = make(map[int64]int)
		m[account] = b
	}
	b[sec]++
	if limit > 0 && b[sec] > limit {
		return true
	}
	return false
}

func (r *AbuseRiskManager) sumSince(b map[int64]int, startSec int64) int {
	if len(b) == 0 {
		return 0
	}
	total := 0
	for sec, v := range b {
		if sec >= startSec {
			total += v
		}
	}
	return total
}
