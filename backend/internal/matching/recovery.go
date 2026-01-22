package matching

import (
	"context"
	"fmt"
	"time"
)

func RecoverSymbol(snapshotStore SnapshotStore, wal WAL, symbol string) (EngineSnapshot, []Trade, error) {
	var base EngineSnapshot
	var afterSeq uint64

	if snapshotStore != nil {
		snap, err := snapshotStore.Get(symbol)
		if err != nil {
			return EngineSnapshot{}, nil, err
		}
		if snap != nil {
			base = *snap
			afterSeq = snap.Seq
		}
	}
	if afterSeq == 0 && base.CapturedAt == 0 && base.State == StatePreOpen {
		base.State = StateContinuous
	}

	ob := NewOrderBook(symbol)
	ob.SetConfig(base.Config)

	for i := range base.OpenOrders {
		o := base.OpenOrders[i].ToOrder()
		if err := ob.AddRestingOrder(o); err != nil {
			return EngineSnapshot{}, nil, err
		}
	}

	events, err := wal.Load(symbol, afterSeq)
	if err != nil {
		return EngineSnapshot{}, nil, err
	}

	var trades []Trade
	for i := range events {
		ev := events[i]
		switch ev.Type {
		case EventSubmit, EventTrigger:
			if ev.Order == nil {
				continue
			}
			o := ev.Order.ToOrder()
			o.Timestamp = ev.Timestamp
			ob.SetClock(func() int64 { return ev.Timestamp })
			ts, err := ob.AddOrder(o)
			ob.ResetClock()
			if err != nil {
				return EngineSnapshot{}, trades, err
			}
			trades = append(trades, ts...)
		case EventCancel:
			if ev.OrderID == "" {
				continue
			}
			ob.CancelOrder(ev.OrderID)
		}
	}

	openOrders := make([]OrderRecord, 0, len(ob.Orders))
	for _, o := range ob.Orders {
		ob.materializeOrder(o)
		openOrders = append(openOrders, o.Record())
	}

	seq := afterSeq
	if len(events) > 0 {
		seq = events[len(events)-1].Seq
	}

	snap := EngineSnapshot{
		Symbol:         symbol,
		CapturedAt:     time.Now().UnixNano(),
		Seq:            seq,
		State:          base.State,
		LastPrice:      base.LastPrice,
		LastPriceTicks: base.LastPriceTicks,
		Config:         ob.cfg,
		OpenOrders:     openOrders,
		Triggers:       base.Triggers,
	}

	return snap, trades, nil
}

func RecoverSymbolFromSnapshotAndEvents(symbol string, snap EngineSnapshot, events []Event) (EngineSnapshot, []Trade, error) {
	wal := NewMemoryWAL()
	for i := range events {
		_ = wal.Append(events[i])
	}
	store := NewMemorySnapshotStore()
	_ = store.Put(symbol, snap)
	return RecoverSymbol(store, wal, symbol)
}

func ValidateSnapshot(symbol string, snap EngineSnapshot) error {
	if snap.Symbol != "" && snap.Symbol != symbol {
		return fmt.Errorf("snapshot symbol mismatch")
	}
	return nil
}

func StartWALPump(ctx context.Context, events <-chan Event, wal WAL) {
	go func() {
		for {
			select {
			case ev := <-events:
				_ = wal.Append(ev)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func StartSnapshotPump(ctx context.Context, engine *Engine, store SnapshotStore, symbol string, every time.Duration) {
	go func() {
		ticker := time.NewTicker(every)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				snap, err := engine.CaptureSnapshot(symbol)
				if err == nil {
					_ = store.Put(symbol, snap)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
