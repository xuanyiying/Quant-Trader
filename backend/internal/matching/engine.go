package matching

import (
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type Engine struct {
	workers map[string]*symbolWorker
	mu      sync.RWMutex
	events  chan Event
}

func NewEngine() *Engine {
	return &Engine{
		workers: make(map[string]*symbolWorker),
		events:  make(chan Event, 65536),
	}
}

func (e *Engine) Events() <-chan Event {
	return e.events
}

func (e *Engine) AddOrderBook(symbol string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.workers[symbol]; !exists {
		w := newSymbolWorker(symbol, e.events)
		e.workers[symbol] = w
		go w.run()
	}
}

func (e *Engine) ConfigureSymbol(symbol string, cfg SymbolConfig) error {
	e.mu.RLock()
	w, exists := e.workers[symbol]
	e.mu.RUnlock()
	if !exists {
		return fmt.Errorf("order book not found for symbol: %s", symbol)
	}

	resp := make(chan configResult, 1)
	w.in <- symbolOp{kind: opConfig, cfg: cfg, configResp: resp}
	r := <-resp
	return r.err
}

func (e *Engine) SetMarketState(symbol string, state MarketState) error {
	e.mu.RLock()
	w, exists := e.workers[symbol]
	e.mu.RUnlock()
	if !exists {
		return fmt.Errorf("order book not found for symbol: %s", symbol)
	}
	resp := make(chan configResult, 1)
	w.in <- symbolOp{kind: opState, state: state, configResp: resp}
	r := <-resp
	return r.err
}

func (e *Engine) SetRiskManager(symbol string, risk RiskManager) error {
	e.mu.RLock()
	w, exists := e.workers[symbol]
	e.mu.RUnlock()
	if !exists {
		return fmt.Errorf("order book not found for symbol: %s", symbol)
	}
	resp := make(chan configResult, 1)
	w.in <- symbolOp{kind: opRisk, risk: risk, configResp: resp}
	r := <-resp
	return r.err
}

func (e *Engine) SetFeeCalculator(symbol string, fee TieredFeeCalculator) error {
	e.mu.RLock()
	w, exists := e.workers[symbol]
	e.mu.RUnlock()
	if !exists {
		return fmt.Errorf("order book not found for symbol: %s", symbol)
	}
	resp := make(chan configResult, 1)
	w.in <- symbolOp{kind: opFee, fee: fee, configResp: resp}
	r := <-resp
	return r.err
}

func (e *Engine) UpdateLastPrice(symbol string, price decimal.Decimal) error {
	e.mu.RLock()
	w, exists := e.workers[symbol]
	e.mu.RUnlock()
	if !exists {
		return fmt.Errorf("order book not found for symbol: %s", symbol)
	}
	resp := make(chan configResult, 1)
	w.in <- symbolOp{kind: opPrice, price: price, configResp: resp}
	r := <-resp
	return r.err
}

func (e *Engine) ProcessOrder(order *Order) ([]Trade, error) {
	e.mu.RLock()
	w, exists := e.workers[order.Symbol]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("order book not found for symbol: %s", order.Symbol)
	}

	resp := make(chan submitResult, 1)
	w.in <- symbolOp{kind: opSubmit, order: order, submitResp: resp}
	r := <-resp
	return r.trades, r.err
}

func (e *Engine) CancelOrder(symbol string, orderID string) (*Order, error) {
	e.mu.RLock()
	w, exists := e.workers[symbol]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("order book not found for symbol: %s", symbol)
	}

	resp := make(chan cancelResult, 1)
	w.in <- symbolOp{kind: opCancel, orderID: orderID, cancelResp: resp}
	r := <-resp
	return r.order, r.err
}

func (e *Engine) GetOrderBook(symbol string) (OrderBookSnapshot, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	w, ok := e.workers[symbol]
	if !ok {
		return OrderBookSnapshot{}, fmt.Errorf("order book not found for symbol: %s", symbol)
	}

	resp := make(chan snapshotResult, 1)
	w.in <- symbolOp{kind: opSnapshot, snapshotResp: resp}
	r := <-resp
	return r.snapshot, nil
}

func (e *Engine) CaptureSnapshot(symbol string) (EngineSnapshot, error) {
	e.mu.RLock()
	w, ok := e.workers[symbol]
	e.mu.RUnlock()
	if !ok {
		return EngineSnapshot{}, fmt.Errorf("order book not found for symbol: %s", symbol)
	}

	resp := make(chan snapshotFullResult, 1)
	w.in <- symbolOp{kind: opFullSnapshot, snapshotFullResp: resp}
	r := <-resp
	return r.snapshot, r.err
}

func (e *Engine) Restore(symbol string, snap EngineSnapshot) error {
	e.mu.RLock()
	w, ok := e.workers[symbol]
	e.mu.RUnlock()
	if !ok {
		return fmt.Errorf("order book not found for symbol: %s", symbol)
	}

	resp := make(chan restoreResult, 1)
	w.in <- symbolOp{kind: opRestore, snapshot: snap, restoreResp: resp}
	r := <-resp
	return r.err
}

type opKind int

const (
	opSubmit opKind = iota
	opCancel
	opSnapshot
	opConfig
	opState
	opRisk
	opFee
	opPrice
	opFullSnapshot
	opRestore
)

type symbolOp struct {
	kind             opKind
	order            *Order
	orderID          string
	cfg              SymbolConfig
	state            MarketState
	risk             RiskManager
	fee              TieredFeeCalculator
	price            decimal.Decimal
	snapshot         EngineSnapshot
	submitResp       chan submitResult
	cancelResp       chan cancelResult
	snapshotResp     chan snapshotResult
	snapshotFullResp chan snapshotFullResult
	restoreResp      chan restoreResult
	configResp       chan configResult
}

type submitResult struct {
	trades []Trade
	err    error
}

type cancelResult struct {
	order *Order
	err   error
}

type snapshotResult struct {
	snapshot OrderBookSnapshot
}

type snapshotFullResult struct {
	snapshot EngineSnapshot
	err      error
}

type restoreResult struct {
	err error
}

type configResult struct {
	err error
}

type symbolWorker struct {
	symbol         string
	ob             *OrderBook
	in             chan symbolOp
	out            chan<- Event
	seq            uint64
	state          MarketState
	lastPrice      decimal.Decimal
	lastPriceTicks int64
	risk           RiskManager
	fee            TieredFeeCalculator
	lastSecond     int64
	ordersInSecond int
	triggers       map[string]*Order
}

func newSymbolWorker(symbol string, out chan<- Event) *symbolWorker {
	return &symbolWorker{
		symbol:         symbol,
		ob:             NewOrderBook(symbol),
		in:             make(chan symbolOp, 4096),
		out:            out,
		state:          StateContinuous,
		lastPrice:      decimal.Zero,
		lastPriceTicks: 0,
		risk:           NoopRiskManager{},
		fee:            TieredFeeCalculator{DefaultRate: decimal.Zero},
		triggers:       make(map[string]*Order),
	}
}

func (w *symbolWorker) run() {
	for op := range w.in {
		switch op.kind {
		case opSubmit:
			ts := time.Now().UnixNano()
			order := op.order
			order.Timestamp = ts

			w.emit(Event{
				Type:      EventSubmit,
				Timestamp: ts,
				Symbol:    w.symbol,
				Order:     ptrOrderRecord(order.Record()),
			})

			if w.state == StateClosed || w.state == StateClosing {
				order.Status = StatusRejected
				w.emit(Event{
					Type:      EventReject,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(order.Record()),
					Reason:    "market is closed",
				})
				op.submitResp <- submitResult{trades: nil, err: fmt.Errorf("market is closed")}
				continue
			}

			if w.applyRateLimit(ts) != nil {
				order.Status = StatusRejected
				w.emit(Event{
					Type:      EventReject,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(order.Record()),
					Reason:    "rate limit exceeded",
				})
				op.submitResp <- submitResult{trades: nil, err: fmt.Errorf("rate limit exceeded")}
				continue
			}

			if order.Type == TypeStop || order.Type == TypeTakeProfit {
				if order.TriggerPrice.LessThanOrEqual(decimal.Zero) {
					order.Status = StatusRejected
					w.emit(Event{
						Type:      EventReject,
						Timestamp: ts,
						Symbol:    w.symbol,
						Order:     ptrOrderRecord(order.Record()),
						Reason:    "invalid trigger price",
					})
					op.submitResp <- submitResult{trades: nil, err: fmt.Errorf("invalid trigger price")}
					continue
				}
				order.Status = StatusOpen
				w.triggers[order.ID] = order
				w.emit(Event{
					Type:      EventOrderUpdate,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(order.Record()),
				})
				op.submitResp <- submitResult{trades: nil, err: nil}
				continue
			}

			if order.Type == TypeIceberg {
				if order.PeakSize.LessThanOrEqual(decimal.Zero) {
					order.Status = StatusRejected
					w.emit(Event{
						Type:      EventReject,
						Timestamp: ts,
						Symbol:    w.symbol,
						Order:     ptrOrderRecord(order.Record()),
						Reason:    "invalid peak size",
					})
					op.submitResp <- submitResult{trades: nil, err: fmt.Errorf("invalid peak size")}
					continue
				}
				totalRemaining := order.Amount.Sub(order.Filled)
				if totalRemaining.LessThanOrEqual(decimal.Zero) {
					order.Status = StatusRejected
					w.emit(Event{
						Type:      EventReject,
						Timestamp: ts,
						Symbol:    w.symbol,
						Order:     ptrOrderRecord(order.Record()),
						Reason:    "invalid remaining",
					})
					op.submitResp <- submitResult{trades: nil, err: fmt.Errorf("invalid remaining")}
					continue
				}
				visible := decimal.Min(order.PeakSize, totalRemaining)
				order.Remaining = visible
				order.HiddenRemaining = totalRemaining.Sub(visible)
			}

			if err := w.applyPriceBand(order); err != nil {
				order.Status = StatusRejected
				w.emit(Event{
					Type:      EventReject,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(order.Record()),
					Reason:    err.Error(),
				})
				op.submitResp <- submitResult{trades: nil, err: err}
				continue
			}

			if err := w.risk.ValidateOrder(w.state, w.lastPrice, w.ob.Snapshot(), order); err != nil {
				order.Status = StatusRejected
				w.emit(Event{
					Type:      EventReject,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(order.Record()),
					Reason:    err.Error(),
				})
				op.submitResp <- submitResult{trades: nil, err: err}
				continue
			}

			w.maybeAlertLargeOrder(ts, order)

			if w.state == StatePreOpen || w.state == StateOpenAuction {
				if err := w.ob.AddRestingOrder(order); err != nil {
					w.emit(Event{
						Type:      EventReject,
						Timestamp: ts,
						Symbol:    w.symbol,
						Order:     ptrOrderRecord(order.Record()),
						Reason:    err.Error(),
					})
					op.submitResp <- submitResult{trades: nil, err: err}
					continue
				}
				w.emit(Event{
					Type:      EventOrderUpdate,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(order.Record()),
				})
				op.submitResp <- submitResult{trades: nil, err: nil}
				continue
			}

			trades, err := w.executeContinuous(order, ts)
			op.submitResp <- submitResult{trades: trades, err: err}
		case opCancel:
			ts := time.Now().UnixNano()
			w.emit(Event{
				Type:      EventCancel,
				Timestamp: ts,
				Symbol:    w.symbol,
				OrderID:   op.orderID,
			})

			if trg, ok := w.triggers[op.orderID]; ok {
				delete(w.triggers, op.orderID)
				trg.Status = StatusCanceled
				w.emit(Event{
					Type:      EventOrderUpdate,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(trg.Record()),
				})
				op.cancelResp <- cancelResult{order: trg, err: nil}
				continue
			}

			order := w.ob.CancelOrder(op.orderID)
			if order == nil {
				w.emit(Event{
					Type:      EventReject,
					Timestamp: ts,
					Symbol:    w.symbol,
					OrderID:   op.orderID,
					Reason:    fmt.Errorf("order not found: %s", op.orderID).Error(),
				})
				op.cancelResp <- cancelResult{order: nil, err: fmt.Errorf("order not found: %s", op.orderID)}
				continue
			}
			w.emit(Event{
				Type:      EventOrderUpdate,
				Timestamp: ts,
				Symbol:    w.symbol,
				Order:     ptrOrderRecord(order.Record()),
			})
			op.cancelResp <- cancelResult{order: order, err: nil}
		case opSnapshot:
			op.snapshotResp <- snapshotResult{snapshot: w.ob.Snapshot()}
		case opFullSnapshot:
			ts := time.Now().UnixNano()
			openOrders := make([]OrderRecord, 0, len(w.ob.Orders))
			for _, o := range w.ob.Orders {
				w.ob.materializeOrder(o)
				openOrders = append(openOrders, o.Record())
			}
			triggers := make([]OrderRecord, 0, len(w.triggers))
			for _, o := range w.triggers {
				w.ob.materializeOrder(o)
				triggers = append(triggers, o.Record())
			}
			op.snapshotFullResp <- snapshotFullResult{
				snapshot: EngineSnapshot{
					Symbol:         w.symbol,
					CapturedAt:     ts,
					Seq:            w.seq,
					State:          w.state,
					LastPrice:      w.lastPrice,
					LastPriceTicks: w.lastPriceTicks,
					Config:         w.ob.cfg,
					OpenOrders:     openOrders,
					Triggers:       triggers,
				},
				err: nil,
			}
		case opRestore:
			w.seq = op.snapshot.Seq
			w.state = op.snapshot.State
			w.lastPrice = op.snapshot.LastPrice
			w.lastPriceTicks = op.snapshot.LastPriceTicks

			w.ob = NewOrderBook(w.symbol)
			w.ob.SetConfig(op.snapshot.Config)
			var restoreErr error
			for i := range op.snapshot.OpenOrders {
				o := op.snapshot.OpenOrders[i].ToOrder()
				if err := w.ob.AddRestingOrder(o); err != nil {
					restoreErr = err
					break
				}
			}
			if restoreErr == nil {
				w.triggers = make(map[string]*Order, len(op.snapshot.Triggers))
				for i := range op.snapshot.Triggers {
					o := op.snapshot.Triggers[i].ToOrder()
					o.Status = StatusOpen
					w.triggers[o.ID] = o
				}
			}
			op.restoreResp <- restoreResult{err: restoreErr}
		case opConfig:
			w.ob.SetConfig(op.cfg)
			op.configResp <- configResult{err: nil}
		case opState:
			old := w.state
			w.state = op.state
			if old == StateOpenAuction && w.state == StateContinuous {
				ts := time.Now().UnixNano()
				w.ob.SetClock(func() int64 { return ts })
				priceTicks, ok := w.ob.CalculateClearingPrice(w.lastPriceTicks)
				var trades []Trade
				if ok {
					trades = w.ob.MatchAuction(priceTicks)
				}
				w.ob.ResetClock()
				w.applyFeesAndEmit(ts, trades)
				w.processTriggers(ts)
			}
			op.configResp <- configResult{err: nil}
		case opRisk:
			if op.risk == nil {
				w.risk = NoopRiskManager{}
			} else {
				w.risk = op.risk
			}
			op.configResp <- configResult{err: nil}
		case opFee:
			w.fee = op.fee
			op.configResp <- configResult{err: nil}
		case opPrice:
			w.setLastPrice(op.price)
			ts := time.Now().UnixNano()
			w.processTriggers(ts)
			op.configResp <- configResult{err: nil}
		}
	}
}

func (w *symbolWorker) emit(ev Event) {
	w.seq++
	ev.Seq = w.seq
	w.out <- ev
	if c, ok := w.risk.(interface{ ConsumeEvent(Event) }); ok {
		c.ConsumeEvent(ev)
	}
}

func (w *symbolWorker) applyRateLimit(ts int64) error {
	limit := w.ob.cfg.MaxOrdersPerSecond
	if limit <= 0 {
		return nil
	}
	sec := ts / 1_000_000_000
	if w.lastSecond != sec {
		w.lastSecond = sec
		w.ordersInSecond = 0
	}
	w.ordersInSecond++
	if w.ordersInSecond > limit {
		return fmt.Errorf("rate limit exceeded")
	}
	return nil
}

func (w *symbolWorker) applyPriceBand(order *Order) error {
	if order.Type == TypeMarket {
		return nil
	}
	if w.ob.cfg.LimitUp.GreaterThan(decimal.Zero) && order.Price.GreaterThan(w.ob.cfg.LimitUp) {
		return fmt.Errorf("price above limit up")
	}
	if w.ob.cfg.LimitDown.GreaterThan(decimal.Zero) && order.Price.LessThan(w.ob.cfg.LimitDown) {
		return fmt.Errorf("price below limit down")
	}
	return nil
}

func (w *symbolWorker) maybeAlertLargeOrder(ts int64, order *Order) {
	if order.Type == TypeMarket {
		return
	}
	threshold := w.ob.cfg.LargeOrderNotional
	if threshold.LessThanOrEqual(decimal.Zero) {
		return
	}
	notional := order.Price.Mul(order.Amount)
	if notional.GreaterThanOrEqual(threshold) {
		w.emit(Event{
			Type:      EventAlert,
			Timestamp: ts,
			Symbol:    w.symbol,
			Order:     ptrOrderRecord(order.Record()),
			Reason:    "large order",
		})
	}
}

func (w *symbolWorker) executeContinuous(order *Order, ts int64) ([]Trade, error) {
	w.ob.SetClock(func() int64 { return ts })
	trades, err := w.ob.AddOrder(order)
	w.ob.ResetClock()
	if err != nil {
		w.emit(Event{
			Type:      EventReject,
			Timestamp: ts,
			Symbol:    w.symbol,
			Order:     ptrOrderRecord(order.Record()),
			Reason:    err.Error(),
		})
		return nil, err
	}

	w.applyFeesAndEmit(ts, trades)
	w.emit(Event{
		Type:      EventOrderUpdate,
		Timestamp: ts,
		Symbol:    w.symbol,
		Order:     ptrOrderRecord(order.Record()),
	})
	w.processTriggers(ts)
	return trades, nil
}

func (w *symbolWorker) applyFeesAndEmit(ts int64, trades []Trade) {
	for i := range trades {
		trades[i].Fee = w.fee.Fee(trades[i].Notional)
		w.setLastPrice(trades[i].Price)
		w.risk.OnTrade(trades[i])
		tr := trades[i]
		w.emit(Event{
			Type:      EventTrade,
			Timestamp: tr.Timestamp,
			Symbol:    w.symbol,
			Trade:     &tr,
		})
	}

	makerUpdates := w.ob.DrainMakerUpdates()
	for i := range makerUpdates {
		u := makerUpdates[i]
		w.emit(Event{
			Type:      EventOrderUpdate,
			Timestamp: ts,
			Symbol:    w.symbol,
			Order:     ptrOrderRecord(u),
		})
	}
}

func (w *symbolWorker) processTriggers(ts int64) {
	if w.state != StateContinuous {
		return
	}
	if w.lastPriceTicks <= 0 {
		return
	}

	for {
		var pending []*Order
		for _, o := range w.triggers {
			if w.shouldTrigger(o) {
				pending = append(pending, o)
			}
		}
		if len(pending) == 0 {
			return
		}

		for i := range pending {
			o := pending[i]
			delete(w.triggers, o.ID)

			o.Type = TypeMarket
			o.Price = decimal.Zero
			o.TriggerPrice = decimal.Zero
			o.Timestamp = ts
			o.Status = StatusNew
			o.Remaining = o.Amount.Sub(o.Filled)

			w.emit(Event{
				Type:      EventTrigger,
				Timestamp: ts,
				Symbol:    w.symbol,
				Order:     ptrOrderRecord(o.Record()),
			})

			w.ob.SetClock(func() int64 { return ts })
			trades, err := w.ob.AddOrder(o)
			w.ob.ResetClock()
			if err != nil {
				w.emit(Event{
					Type:      EventReject,
					Timestamp: ts,
					Symbol:    w.symbol,
					Order:     ptrOrderRecord(o.Record()),
					Reason:    err.Error(),
				})
				continue
			}
			w.applyFeesAndEmit(ts, trades)
			w.emit(Event{
				Type:      EventOrderUpdate,
				Timestamp: ts,
				Symbol:    w.symbol,
				Order:     ptrOrderRecord(o.Record()),
			})
		}
	}
}

func (w *symbolWorker) shouldTrigger(order *Order) bool {
	if order.Type == TypeStop {
		if order.Side == SideBuy {
			return w.lastPriceTicks >= order.TriggerTicks
		}
		return w.lastPriceTicks <= order.TriggerTicks
	}
	if order.Type == TypeTakeProfit {
		if order.Side == SideBuy {
			return w.lastPriceTicks <= order.TriggerTicks
		}
		return w.lastPriceTicks >= order.TriggerTicks
	}
	return false
}

func (w *symbolWorker) setLastPrice(price decimal.Decimal) {
	w.lastPrice = price
	if w.ob.cfg.TickSize.GreaterThan(decimal.Zero) {
		q := price.Div(w.ob.cfg.TickSize).Round(0)
		w.lastPriceTicks = q.IntPart()
	}
}

func ptrOrderRecord(r OrderRecord) *OrderRecord {
	return &r
}

func Replay(symbol string, cfg SymbolConfig, events []Event) (OrderBookSnapshot, []Trade, error) {
	ob := NewOrderBook(symbol)
	ob.SetConfig(cfg)

	var produced []Trade
	var recorded []Trade

	for _, ev := range events {
		switch ev.Type {
		case EventSubmit, EventTrigger:
			if ev.Order == nil {
				continue
			}
			o := ev.Order.ToOrder()
			if o.Type == TypeStop || o.Type == TypeTakeProfit {
				continue
			}
			o.Timestamp = ev.Timestamp
			ob.SetClock(func() int64 { return ev.Timestamp })
			trades, err := ob.AddOrder(o)
			ob.ResetClock()
			if err != nil {
				return OrderBookSnapshot{}, produced, err
			}
			produced = append(produced, trades...)
		case EventCancel:
			if ev.OrderID == "" {
				continue
			}
			ob.CancelOrder(ev.OrderID)
		case EventTrade:
			if ev.Trade != nil {
				recorded = append(recorded, *ev.Trade)
			}
		}
	}

	if len(recorded) > 0 {
		if !tradesEqual(produced, recorded) {
			return ob.Snapshot(), produced, fmt.Errorf("replay trade mismatch")
		}
	}

	return ob.Snapshot(), produced, nil
}

func tradesEqual(a []Trade, b []Trade) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].MakerByID != b[i].MakerByID {
			return false
		}
		if a[i].TakerByID != b[i].TakerByID {
			return false
		}
		if !a[i].Price.Equal(b[i].Price) {
			return false
		}
		if !a[i].Amount.Equal(b[i].Amount) {
			return false
		}
	}
	return true
}
