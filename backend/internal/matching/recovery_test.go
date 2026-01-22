package matching

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRecoverSymbol_FromSnapshotAndWAL(t *testing.T) {
	symbol := "REC-USDT"
	cfg := SymbolConfig{
		TickSize: decimal.RequireFromString("0.00000001"),
		LotSize:  decimal.RequireFromString("0.00000001"),
	}

	ob := NewOrderBook(symbol)
	ob.SetConfig(cfg)

	sell := NewOrder("s1", symbol, SideSell, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	assert.NoError(t, ob.AddRestingOrder(sell))

	open := make([]OrderRecord, 0, len(ob.Orders))
	for _, o := range ob.Orders {
		ob.materializeOrder(o)
		open = append(open, o.Record())
	}

	snap := EngineSnapshot{
		Symbol:         symbol,
		Seq:            100,
		State:          StateContinuous,
		LastPrice:      decimal.Zero,
		LastPriceTicks: 0,
		Config:         cfg,
		OpenOrders:     open,
	}

	wal := NewMemoryWAL()
	buy := NewOrder("b1", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	assert.NoError(t, ob.validateOrder(buy))
	ev := Event{
		Seq:       101,
		Type:      EventSubmit,
		Timestamp: 123,
		Symbol:    symbol,
		Order:     ptrOrderRecord(buy.Record()),
	}
	assert.NoError(t, wal.Append(ev))

	store := NewMemorySnapshotStore()
	assert.NoError(t, store.Put(symbol, snap))

	outSnap, trades, err := RecoverSymbol(store, wal, symbol)
	assert.NoError(t, err)
	assert.Len(t, trades, 1)
	assert.Equal(t, "s1", trades[0].MakerByID)
	assert.Equal(t, "b1", trades[0].TakerByID)

	assert.Equal(t, uint64(101), outSnap.Seq)
	assert.Len(t, outSnap.OpenOrders, 0)
}
