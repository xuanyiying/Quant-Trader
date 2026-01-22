package matching

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMatchingEngine_LimitOrder(t *testing.T) {
	engine := NewEngine()
	symbol := "BTC-USDT"
	engine.AddOrderBook(symbol)

	// 1. Place Sell Limit Order: 100 @ 100.0
	sellOrder := NewOrder("order-1", symbol, SideSell, TypeLimit, decimal.NewFromFloat(100.0), decimal.NewFromFloat(1.0))
	trades, err := engine.ProcessOrder(sellOrder)
	assert.NoError(t, err)
	assert.Empty(t, trades)

	// Check Book
	snap, err := engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.Nil(t, snap.BestBid)
	assert.NotNil(t, snap.BestAsk)
	assert.Equal(t, decimal.NewFromFloat(100.0).String(), snap.BestAsk.Price.String())

	// 2. Place Buy Limit Order: 50 @ 100.0 (Partial Match)
	buyOrder := NewOrder("order-2", symbol, SideBuy, TypeLimit, decimal.NewFromFloat(100.0), decimal.NewFromFloat(0.5))
	trades, err = engine.ProcessOrder(buyOrder)
	assert.NoError(t, err)
	assert.Len(t, trades, 1)

	trade := trades[0]
	assert.Equal(t, "order-1", trade.MakerByID)
	assert.Equal(t, "order-2", trade.TakerByID)
	assert.Equal(t, decimal.NewFromFloat(100.0).String(), trade.Price.String())
	assert.Equal(t, decimal.NewFromFloat(0.5).String(), trade.Amount.String())

	// Check Order Status
	assert.Equal(t, StatusPartiallyFilled, sellOrder.Status)
	assert.Equal(t, StatusFilled, buyOrder.Status)
	assert.Equal(t, decimal.NewFromFloat(0.5).String(), sellOrder.Remaining.String())

	// 3. Place Buy Limit Order: 60 @ 101.0 (Match remaining 50 @ 100.0, rest on book)
	// Note: Taker buys at Maker's price (100.0). The remaining 10 should sit at 101.0
	buyOrder2 := NewOrder("order-3", symbol, SideBuy, TypeLimit, decimal.NewFromFloat(101.0), decimal.NewFromFloat(0.6))
	trades, err = engine.ProcessOrder(buyOrder2)
	assert.NoError(t, err)
	assert.Len(t, trades, 1)

	trade2 := trades[0]
	assert.Equal(t, "order-1", trade2.MakerByID)
	assert.Equal(t, decimal.NewFromFloat(100.0).String(), trade2.Price.String())
	assert.Equal(t, decimal.NewFromFloat(0.5).String(), trade2.Amount.String())

	// Check Sell Order Filled
	assert.Equal(t, StatusFilled, sellOrder.Status)
	assert.True(t, sellOrder.Remaining.IsZero())

	// Check Buy Order 2 Remaining on Book
	assert.Equal(t, StatusPartiallyFilled, buyOrder2.Status)
	assert.Equal(t, decimal.NewFromFloat(0.1).String(), buyOrder2.Remaining.String())

	// Asks should be empty now
	snap, err = engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.Nil(t, snap.BestAsk)
	// Bids should have 0.1 @ 101.0
	assert.NotNil(t, snap.BestBid)
	assert.Equal(t, decimal.NewFromFloat(101.0).String(), snap.BestBid.Price.String())
}

func TestMatchingEngine_PriceTimePriority(t *testing.T) {
	engine := NewEngine()
	symbol := "ETH-USDT"
	engine.AddOrderBook(symbol)

	// Place two Sell Orders at same price
	sell1 := NewOrder("s1", symbol, SideSell, TypeLimit, decimal.NewFromFloat(2000), decimal.NewFromFloat(1))
	sell2 := NewOrder("s2", symbol, SideSell, TypeLimit, decimal.NewFromFloat(2000), decimal.NewFromFloat(1))

	engine.ProcessOrder(sell1)
	engine.ProcessOrder(sell2)

	// Place Buy for 1.5
	buy := NewOrder("b1", symbol, SideBuy, TypeLimit, decimal.NewFromFloat(2000), decimal.NewFromFloat(1.5))
	trades, _ := engine.ProcessOrder(buy)

	assert.Len(t, trades, 2)
	// Should match s1 first (Time priority)
	assert.Equal(t, "s1", trades[0].MakerByID)
	assert.Equal(t, decimal.NewFromFloat(1).String(), trades[0].Amount.String())

	// Then s2
	assert.Equal(t, "s2", trades[1].MakerByID)
	assert.Equal(t, decimal.NewFromFloat(0.5).String(), trades[1].Amount.String())
}

func TestSkipList_RemoveLevel(t *testing.T) {
	sl := NewSkipList(true) // Descending

	priceTicks := int64(100)
	o1 := NewOrder("1", "S", SideBuy, TypeLimit, decimal.NewFromInt(100), decimal.NewFromInt(1))
	o1.PriceTicks = priceTicks
	o1.AmountLots = 1
	o1.RemainingLots = 1

	sl.Insert(priceTicks, o1)

	assert.NotNil(t, sl.Get(priceTicks))

	sl.RemoveLevel(priceTicks)

	assert.Nil(t, sl.Get(priceTicks))
}

func TestMatchingEngine_MarketOrder(t *testing.T) {
	engine := NewEngine()
	symbol := "SOL-USDT"
	engine.AddOrderBook(symbol)

	// Setup liquidity: Sell 10 @ 20, Sell 10 @ 21
	engine.ProcessOrder(NewOrder("s1", symbol, SideSell, TypeLimit, decimal.NewFromFloat(20), decimal.NewFromFloat(10)))
	engine.ProcessOrder(NewOrder("s2", symbol, SideSell, TypeLimit, decimal.NewFromFloat(21), decimal.NewFromFloat(10)))

	// Market Buy 15
	marketBuy := NewOrder("b1", symbol, SideBuy, TypeMarket, decimal.Zero, decimal.NewFromFloat(15))
	trades, _ := engine.ProcessOrder(marketBuy)

	assert.Len(t, trades, 2)
	assert.Equal(t, "s1", trades[0].MakerByID)
	assert.Equal(t, decimal.NewFromFloat(10).String(), trades[0].Amount.String())
	assert.Equal(t, decimal.NewFromFloat(20).String(), trades[0].Price.String())

	assert.Equal(t, "s2", trades[1].MakerByID)
	assert.Equal(t, decimal.NewFromFloat(5).String(), trades[1].Amount.String())
	assert.Equal(t, decimal.NewFromFloat(21).String(), trades[1].Price.String())
}

func TestOrderQueue_TotalVolumeUpdatesOnPartialFill(t *testing.T) {
	engine := NewEngine()
	symbol := "XRP-USDT"
	engine.AddOrderBook(symbol)

	engine.ProcessOrder(NewOrder("s1", symbol, SideSell, TypeLimit, decimal.NewFromInt(10), decimal.NewFromInt(1)))

	snap, err := engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.NotNil(t, snap.BestAsk)
	assert.Equal(t, decimal.NewFromInt(1).String(), snap.BestAsk.Volume.String())

	engine.ProcessOrder(NewOrder("b1", symbol, SideBuy, TypeLimit, decimal.NewFromInt(10), decimal.NewFromFloat(0.4)))
	snap, err = engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.NotNil(t, snap.BestAsk)
	assert.Equal(t, decimal.NewFromFloat(0.6).String(), snap.BestAsk.Volume.String())
}

func TestSymbolConfig_TickLotMinValidation(t *testing.T) {
	engine := NewEngine()
	symbol := "CFG-USDT"
	engine.AddOrderBook(symbol)

	cfg := SymbolConfig{
		TickSize:    decimal.RequireFromString("0.1"),
		LotSize:     decimal.RequireFromString("0.01"),
		MinQty:      decimal.RequireFromString("0.01"),
		MinNotional: decimal.RequireFromString("1"),
	}
	assert.NoError(t, engine.ConfigureSymbol(symbol, cfg))

	o1 := NewOrder("o1", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10.05"), decimal.RequireFromString("0.01"))
	_, err := engine.ProcessOrder(o1)
	assert.Error(t, err)
	assert.Equal(t, StatusRejected, o1.Status)

	o2 := NewOrder("o2", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10.0"), decimal.RequireFromString("0.005"))
	_, err = engine.ProcessOrder(o2)
	assert.Error(t, err)
	assert.Equal(t, StatusRejected, o2.Status)

	o3 := NewOrder("o3", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10.0"), decimal.RequireFromString("0.1"))
	_, err = engine.ProcessOrder(o3)
	assert.NoError(t, err)
}

func TestAuditEvents_ReplayMatchesTrades(t *testing.T) {
	engine := NewEngine()
	symbol := "RPL-USDT"
	engine.AddOrderBook(symbol)

	_, err := engine.ProcessOrder(NewOrder("s1", symbol, SideSell, TypeLimit, decimal.NewFromInt(10), decimal.NewFromInt(1)))
	assert.NoError(t, err)
	_, err = engine.ProcessOrder(NewOrder("b1", symbol, SideBuy, TypeLimit, decimal.NewFromInt(10), decimal.NewFromInt(1)))
	assert.NoError(t, err)

	var events []Event
	for i := 0; i < 6; i++ {
		events = append(events, <-engine.Events())
	}

	_, produced, err := Replay(symbol, SymbolConfig{}, events)
	assert.NoError(t, err)

	tradeEvents := 0
	for i := range events {
		if events[i].Type == EventTrade {
			tradeEvents++
		}
	}

	assert.Equal(t, tradeEvents, len(produced))
}

func TestIcebergOrder_Replenish(t *testing.T) {
	engine := NewEngine()
	symbol := "ICE-USDT"
	engine.AddOrderBook(symbol)

	sell := NewOrder("ice-s1", symbol, SideSell, TypeIceberg, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	sell.PeakSize = decimal.RequireFromString("0.3")
	_, err := engine.ProcessOrder(sell)
	assert.NoError(t, err)

	snap, err := engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.NotNil(t, snap.BestAsk)
	assert.Equal(t, decimal.RequireFromString("0.3").String(), snap.BestAsk.Volume.String())

	buy := NewOrder("ice-b1", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("0.5"))
	trades, err := engine.ProcessOrder(buy)
	assert.NoError(t, err)
	assert.Len(t, trades, 2)

	snap, err = engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.NotNil(t, snap.BestAsk)
	assert.Equal(t, decimal.RequireFromString("0.1").String(), snap.BestAsk.Volume.String())
}

func TestStopOrder_TriggersOnPriceUpdate(t *testing.T) {
	engine := NewEngine()
	symbol := "STP-USDT"
	engine.AddOrderBook(symbol)

	stopBuy := NewOrder("stp-b1", symbol, SideBuy, TypeStop, decimal.Zero, decimal.RequireFromString("1"))
	stopBuy.TriggerPrice = decimal.RequireFromString("10")
	_, err := engine.ProcessOrder(stopBuy)
	assert.NoError(t, err)
	assert.Equal(t, StatusOpen, stopBuy.Status)

	sell := NewOrder("stp-s1", symbol, SideSell, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	_, err = engine.ProcessOrder(sell)
	assert.NoError(t, err)

	assert.NoError(t, engine.UpdateLastPrice(symbol, decimal.RequireFromString("10")))
	assert.Equal(t, StatusFilled, stopBuy.Status)

	snap, err := engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.Nil(t, snap.BestAsk)
	assert.Nil(t, snap.BestBid)
}

func TestOpenAuction_MatchesOnTransitionToContinuous(t *testing.T) {
	engine := NewEngine()
	symbol := "AUC-USDT"
	engine.AddOrderBook(symbol)

	assert.NoError(t, engine.SetMarketState(symbol, StateOpenAuction))

	buy := NewOrder("auc-b1", symbol, SideBuy, TypeLimit, decimal.RequireFromString("10"), decimal.RequireFromString("1"))
	sell := NewOrder("auc-s1", symbol, SideSell, TypeLimit, decimal.RequireFromString("9"), decimal.RequireFromString("1"))
	_, err := engine.ProcessOrder(buy)
	assert.NoError(t, err)
	_, err = engine.ProcessOrder(sell)
	assert.NoError(t, err)

	snap, err := engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.NotNil(t, snap.BestBid)
	assert.NotNil(t, snap.BestAsk)

	assert.NoError(t, engine.SetMarketState(symbol, StateContinuous))

	snap, err = engine.GetOrderBook(symbol)
	assert.NoError(t, err)
	assert.Nil(t, snap.BestBid)
	assert.Nil(t, snap.BestAsk)
}
