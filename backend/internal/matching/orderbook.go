package matching

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

type OrderBook struct {
	Symbol string
	Bids   *SkipList // Descending (High -> Low)
	Asks   *SkipList // Ascending (Low -> High)
	// OrderID -> Order mapping for O(1) cancellation/lookup
	Orders       map[string]*Order
	seq          uint64
	now          func() int64
	defaultNow   func() int64
	cfg          SymbolConfig
	makerUpdates []OrderRecord
}

type LevelSnapshot struct {
	Price      decimal.Decimal `json:"price"`
	Volume     decimal.Decimal `json:"volume"`
	OrderCount int             `json:"order_count"`
}

type OrderBookSnapshot struct {
	Symbol  string         `json:"symbol"`
	BestBid *LevelSnapshot `json:"best_bid,omitempty"`
	BestAsk *LevelSnapshot `json:"best_ask,omitempty"`
}

func (ob *OrderBook) Snapshot() OrderBookSnapshot {
	snap := OrderBookSnapshot{Symbol: ob.Symbol}

	if bid := ob.Bids.Best(); bid != nil {
		snap.BestBid = &LevelSnapshot{
			Price:      ob.priceFromTicks(bid.PriceTicks),
			Volume:     ob.qtyFromLots(bid.TotalLots),
			OrderCount: bid.Len(),
		}
	}

	if ask := ob.Asks.Best(); ask != nil {
		snap.BestAsk = &LevelSnapshot{
			Price:      ob.priceFromTicks(ask.PriceTicks),
			Volume:     ob.qtyFromLots(ask.TotalLots),
			OrderCount: ask.Len(),
		}
	}

	return snap
}

func NewOrderBook(symbol string) *OrderBook {
	defaultNow := func() int64 { return time.Now().UnixNano() }
	return &OrderBook{
		Symbol:       symbol,
		Bids:         NewSkipList(true),
		Asks:         NewSkipList(false),
		Orders:       make(map[string]*Order),
		now:          defaultNow,
		defaultNow:   defaultNow,
		makerUpdates: make([]OrderRecord, 0, 16),
		cfg: SymbolConfig{
			TickSize: decimal.RequireFromString("0.00000001"),
			LotSize:  decimal.RequireFromString("0.00000001"),
		},
	}
}

func (ob *OrderBook) SetConfig(cfg SymbolConfig) {
	if cfg.TickSize.GreaterThan(decimal.Zero) {
		ob.cfg.TickSize = cfg.TickSize
	}
	if cfg.LotSize.GreaterThan(decimal.Zero) {
		ob.cfg.LotSize = cfg.LotSize
	}
	if cfg.MinQty.GreaterThan(decimal.Zero) {
		ob.cfg.MinQty = cfg.MinQty
	}
	if cfg.MinNotional.GreaterThan(decimal.Zero) {
		ob.cfg.MinNotional = cfg.MinNotional
	}
	if cfg.LimitUp.GreaterThan(decimal.Zero) {
		ob.cfg.LimitUp = cfg.LimitUp
	}
	if cfg.LimitDown.GreaterThan(decimal.Zero) {
		ob.cfg.LimitDown = cfg.LimitDown
	}
	if cfg.LargeOrderNotional.GreaterThan(decimal.Zero) {
		ob.cfg.LargeOrderNotional = cfg.LargeOrderNotional
	}
	if cfg.MaxOrdersPerSecond > 0 {
		ob.cfg.MaxOrdersPerSecond = cfg.MaxOrdersPerSecond
	}
}

func (ob *OrderBook) SetClock(now func() int64) {
	ob.now = now
}

func (ob *OrderBook) ResetClock() {
	ob.now = ob.defaultNow
}

func (ob *OrderBook) priceFromTicks(ticks int64) decimal.Decimal {
	return ob.cfg.TickSize.Mul(decimal.NewFromInt(ticks))
}

func (ob *OrderBook) qtyFromLots(lots int64) decimal.Decimal {
	return ob.cfg.LotSize.Mul(decimal.NewFromInt(lots))
}

func (ob *OrderBook) quantizePrice(price decimal.Decimal) (int64, error) {
	if ob.cfg.TickSize.LessThanOrEqual(decimal.Zero) {
		return 0, fmt.Errorf("tick size must be > 0")
	}
	q := price.Div(ob.cfg.TickSize)
	if !q.Equal(q.Truncate(0)) {
		return 0, fmt.Errorf("price not aligned to tick size")
	}
	return q.IntPart(), nil
}

func (ob *OrderBook) quantizeQty(qty decimal.Decimal) (int64, error) {
	if ob.cfg.LotSize.LessThanOrEqual(decimal.Zero) {
		return 0, fmt.Errorf("lot size must be > 0")
	}
	q := qty.Div(ob.cfg.LotSize)
	if !q.Equal(q.Truncate(0)) {
		return 0, fmt.Errorf("amount not aligned to lot size")
	}
	return q.IntPart(), nil
}

func (ob *OrderBook) materializeOrder(o *Order) {
	if o.Type != TypeMarket {
		o.Price = ob.priceFromTicks(o.PriceTicks)
	}
	o.Amount = ob.qtyFromLots(o.AmountLots)
	o.Filled = ob.qtyFromLots(o.FilledLots)
	o.Remaining = ob.qtyFromLots(o.RemainingLots)

	if o.TriggerTicks != 0 {
		o.TriggerPrice = ob.priceFromTicks(o.TriggerTicks)
	}
	if o.PeakLots != 0 {
		o.PeakSize = ob.qtyFromLots(o.PeakLots)
	}
	if o.HiddenRemainingLots != 0 {
		o.HiddenRemaining = ob.qtyFromLots(o.HiddenRemainingLots)
	}
}

func (ob *OrderBook) validateOrder(order *Order) error {
	if order.Amount.LessThanOrEqual(decimal.Zero) && order.AmountLots == 0 {
		return fmt.Errorf("amount must be > 0")
	}

	if order.Filled.IsNegative() || order.FilledLots < 0 {
		return fmt.Errorf("filled must be >= 0")
	}

	if order.AmountLots == 0 {
		lots, err := ob.quantizeQty(order.Amount)
		if err != nil {
			return err
		}
		order.AmountLots = lots
	}
	if order.FilledLots == 0 && order.Filled.GreaterThan(decimal.Zero) {
		lots, err := ob.quantizeQty(order.Filled)
		if err != nil {
			return err
		}
		order.FilledLots = lots
	}

	if order.AmountLots <= 0 {
		return fmt.Errorf("amount must be > 0")
	}
	if order.FilledLots < 0 {
		return fmt.Errorf("filled must be >= 0")
	}
	if order.FilledLots > order.AmountLots {
		return fmt.Errorf("filled exceeds amount")
	}

	totalRemainingLots := order.AmountLots - order.FilledLots

	if order.Type != TypeMarket && order.PriceTicks == 0 {
		if order.Price.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("price must be > 0")
		}
		ticks, err := ob.quantizePrice(order.Price)
		if err != nil {
			return err
		}
		order.PriceTicks = ticks
	}

	if order.Type == TypeStop || order.Type == TypeTakeProfit {
		if order.TriggerTicks == 0 {
			if order.TriggerPrice.LessThanOrEqual(decimal.Zero) {
				return fmt.Errorf("invalid trigger price")
			}
			ticks, err := ob.quantizePrice(order.TriggerPrice)
			if err != nil {
				return err
			}
			order.TriggerTicks = ticks
		}
	}

	if order.Type == TypeIceberg {
		if order.PeakLots == 0 {
			if order.PeakSize.LessThanOrEqual(decimal.Zero) {
				return fmt.Errorf("invalid peak size")
			}
			peakLots, err := ob.quantizeQty(order.PeakSize)
			if err != nil {
				return err
			}
			order.PeakLots = peakLots
		}
		if order.PeakLots <= 0 {
			return fmt.Errorf("invalid peak size")
		}

		if order.RemainingLots == 0 && order.HiddenRemainingLots == 0 {
			visibleLots := minInt64(order.PeakLots, totalRemainingLots)
			order.RemainingLots = visibleLots
			order.HiddenRemainingLots = totalRemainingLots - visibleLots
		} else {
			if order.RemainingLots+order.HiddenRemainingLots != totalRemainingLots {
				return fmt.Errorf("iceberg remaining mismatch")
			}
		}
	} else {
		order.RemainingLots = totalRemainingLots
	}

	if ob.cfg.MinQty.GreaterThan(decimal.Zero) {
		minLots, err := ob.quantizeQty(ob.cfg.MinQty)
		if err == nil && order.AmountLots < minLots {
			return fmt.Errorf("amount below min qty")
		}
	}

	if ob.cfg.MinNotional.GreaterThan(decimal.Zero) && order.Type != TypeMarket {
		notional := ob.priceFromTicks(order.PriceTicks).Mul(ob.qtyFromLots(order.AmountLots))
		if notional.LessThan(ob.cfg.MinNotional) {
			return fmt.Errorf("notional below min notional")
		}
	}

	ob.materializeOrder(order)
	return nil
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// AddOrder processes an incoming order and returns generated trades
func (ob *OrderBook) AddOrder(order *Order) ([]Trade, error) {
	if order.Symbol != ob.Symbol {
		return nil, fmt.Errorf("order symbol mismatch")
	}

	ob.makerUpdates = ob.makerUpdates[:0]

	if err := ob.validateOrder(order); err != nil {
		order.Status = StatusRejected
		return nil, err
	}

	if order.Type == TypeFOK {
		if !ob.canFullyFill(order) {
			order.Status = StatusRejected
			return nil, fmt.Errorf("fok not fillable")
		}
	}

	var trades []Trade
	switch order.Type {
	case TypeMarket:
		trades = ob.matchMarketOrder(order)
	default:
		trades = ob.matchLimitOrder(order)
	}

	if order.RemainingLots == 0 {
		order.Status = StatusFilled
		ob.materializeOrder(order)
		return trades, nil
	}

	if order.FilledLots > 0 {
		order.Status = StatusPartiallyFilled
	} else {
		order.Status = StatusNew
	}

	switch order.Type {
	case TypeLimit, TypeIceberg:
		if order.RemainingLots > 0 {
			ob.addOrderToBook(order)
			if order.FilledLots == 0 {
				order.Status = StatusOpen
			}
		}
	case TypeIOC, TypeMarket:
		if order.FilledLots > 0 {
			order.Status = StatusPartiallyFilledCanceled
		} else {
			order.Status = StatusCanceled
		}
	}

	ob.materializeOrder(order)
	return trades, nil
}

func (ob *OrderBook) AddRestingOrder(order *Order) error {
	if order.Symbol != ob.Symbol {
		return fmt.Errorf("order symbol mismatch")
	}
	ob.makerUpdates = ob.makerUpdates[:0]

	if order.Type == TypeMarket || order.Type == TypeIOC || order.Type == TypeFOK {
		order.Status = StatusRejected
		return fmt.Errorf("order type not allowed in resting mode")
	}

	if err := ob.validateOrder(order); err != nil {
		order.Status = StatusRejected
		return err
	}

	if order.RemainingLots <= 0 {
		order.Status = StatusRejected
		return fmt.Errorf("remaining must be > 0")
	}

	ob.addOrderToBook(order)
	order.Status = StatusOpen
	ob.materializeOrder(order)
	return nil
}

func (ob *OrderBook) DrainMakerUpdates() []OrderRecord {
	if len(ob.makerUpdates) == 0 {
		return nil
	}
	out := make([]OrderRecord, len(ob.makerUpdates))
	copy(out, ob.makerUpdates)
	ob.makerUpdates = ob.makerUpdates[:0]
	return out
}

func (ob *OrderBook) CancelOrder(orderID string) *Order {
	order, ok := ob.Orders[orderID]
	if !ok {
		return nil
	}

	// Remove from SkipList
	if order.Side == SideBuy {
		ob.Bids.Remove(order.PriceTicks, orderID)
	} else {
		ob.Asks.Remove(order.PriceTicks, orderID)
	}

	delete(ob.Orders, orderID)
	order.Status = StatusCanceled
	ob.materializeOrder(order)
	return order
}

func (ob *OrderBook) addOrderToBook(order *Order) {
	if order.Side == SideBuy {
		ob.Bids.Insert(order.PriceTicks, order)
	} else {
		ob.Asks.Insert(order.PriceTicks, order)
	}
	ob.Orders[order.ID] = order
}

func (ob *OrderBook) nextTradeID() string {
	ob.seq++
	return fmt.Sprintf("%s-%d", ob.Symbol, ob.seq)
}

func (ob *OrderBook) canFullyFill(order *Order) bool {
	if order.RemainingLots <= 0 {
		return true
	}

	need := order.RemainingLots
	var fill int64

	if order.Side == SideBuy {
		ob.Asks.ForEachLevel0(func(q *OrderQueue) bool {
			if order.Type != TypeMarket && q.PriceTicks > order.PriceTicks {
				return false
			}
			fill += q.TotalLots
			return fill < need
		})
		return fill >= need
	}

	ob.Bids.ForEachLevel0(func(q *OrderQueue) bool {
		if order.Type != TypeMarket && q.PriceTicks < order.PriceTicks {
			return false
		}
		fill += q.TotalLots
		return fill < need
	})
	return fill >= need
}

func (ob *OrderBook) matchLimitOrder(order *Order) []Trade {
	var trades []Trade
	var oppositeBook *SkipList

	if order.Side == SideBuy {
		oppositeBook = ob.Asks
	} else {
		oppositeBook = ob.Bids
	}

	for order.RemainingLots > 0 {
		bestQueue := oppositeBook.Best()
		if bestQueue == nil {
			break
		}

		// Price Check
		// Buy: Limit Price >= Best Ask
		// Sell: Limit Price <= Best Bid
		if order.Side == SideBuy {
			if order.PriceTicks < bestQueue.PriceTicks {
				break
			}
		} else {
			if order.PriceTicks > bestQueue.PriceTicks {
				break
			}
		}

		// Match against the queue
		newTrades := ob.matchWithQueue(order, bestQueue)
		trades = append(trades, newTrades...)

		// If queue is empty, remove it (implicitly handled by SkipList.Remove if we were removing one by one,
		// but here we might need to cleanup empty queues from the head).
		// Our SkipList.Best() just peeks.
		// matchWithQueue modifies the orders in the queue.
		// If the head order of the queue is filled, it should be removed.
		// The `matchWithQueue` helper will handle popping filled orders.

		// Clean up empty queue/level if needed
		if bestQueue.Len() == 0 {
			oppositeBook.RemoveLevel(bestQueue.PriceTicks)
		}
	}
	return trades
}

func (ob *OrderBook) matchMarketOrder(order *Order) []Trade {
	var trades []Trade
	var oppositeBook *SkipList

	if order.Side == SideBuy {
		oppositeBook = ob.Asks
	} else {
		oppositeBook = ob.Bids
	}

	for order.RemainingLots > 0 {
		bestQueue := oppositeBook.Best()
		if bestQueue == nil {
			break
		}

		// No price check for Market Order (except maybe bounds protection)
		newTrades := ob.matchWithQueue(order, bestQueue)
		trades = append(trades, newTrades...)

		if bestQueue.Len() == 0 {
			oppositeBook.RemoveLevel(bestQueue.PriceTicks)
		}
	}
	return trades
}

func (ob *OrderBook) matchWithQueue(takerOrder *Order, makerQueue *OrderQueue) []Trade {
	var trades []Trade

	for takerOrder.RemainingLots > 0 && makerQueue.Len() > 0 {
		makerOrder := makerQueue.Head()

		matchLots := minInt64(takerOrder.RemainingLots, makerOrder.RemainingLots)
		price := ob.priceFromTicks(makerOrder.PriceTicks)
		amount := ob.qtyFromLots(matchLots)

		// Create Trade
		trade := Trade{
			ID:             ob.nextTradeID(),
			Symbol:         ob.Symbol,
			MakerByID:      makerOrder.ID,
			TakerByID:      takerOrder.ID,
			MakerAccountID: makerOrder.AccountID,
			TakerAccountID: takerOrder.AccountID,
			Price:          price,
			Amount:         amount,
			Notional:       price.Mul(amount),
			Fee:            decimal.Zero,
			Side:           takerOrder.Side.String(),
			Timestamp:      ob.now(),
			IsBuyerMake:    makerOrder.Side == SideBuy,
		}
		trades = append(trades, trade)

		takerOrder.FilledLots += matchLots
		takerOrder.RemainingLots -= matchLots

		makerQueue.TotalLots -= matchLots
		if makerQueue.TotalLots < 0 {
			makerQueue.TotalLots = 0
		}

		makerOrder.FilledLots += matchLots
		makerOrder.RemainingLots -= matchLots

		if makerOrder.RemainingLots == 0 {
			makerQueue.Pop()
			if makerOrder.Type == TypeIceberg && makerOrder.HiddenRemainingLots > 0 && makerOrder.PeakLots > 0 {
				replenishLots := minInt64(makerOrder.PeakLots, makerOrder.HiddenRemainingLots)
				makerOrder.HiddenRemainingLots -= replenishLots
				makerOrder.RemainingLots = replenishLots
				makerOrder.Status = StatusPartiallyFilled
				makerQueue.Add(makerOrder)
				ob.materializeOrder(makerOrder)
				ob.makerUpdates = append(ob.makerUpdates, makerOrder.Record())
				continue
			}

			makerOrder.Status = StatusFilled
			ob.materializeOrder(makerOrder)
			ob.makerUpdates = append(ob.makerUpdates, makerOrder.Record())
			delete(ob.Orders, makerOrder.ID)

		} else {
			makerOrder.Status = StatusPartiallyFilled
			ob.materializeOrder(makerOrder)
			ob.makerUpdates = append(ob.makerUpdates, makerOrder.Record())
		}
	}

	return trades
}

func (ob *OrderBook) CalculateClearingPrice(lastPriceTicks int64) (int64, bool) {
	priceSet := make(map[int64]struct{})
	ob.Bids.ForEachLevel0(func(q *OrderQueue) bool {
		priceSet[q.PriceTicks] = struct{}{}
		return true
	})
	ob.Asks.ForEachLevel0(func(q *OrderQueue) bool {
		priceSet[q.PriceTicks] = struct{}{}
		return true
	})

	if len(priceSet) == 0 {
		return 0, false
	}

	prices := make([]int64, 0, len(priceSet))
	for p := range priceSet {
		prices = append(prices, p)
	}
	sort.Slice(prices, func(i, j int) bool { return prices[i] < prices[j] })

	bestPrice := prices[0]
	bestExec := int64(0)
	bestImbalance := int64(0)
	bestDistance := int64(0)
	hasBest := false

	for i := range prices {
		p := prices[i]

		var buyLots int64
		ob.Bids.ForEachLevel0(func(q *OrderQueue) bool {
			if q.PriceTicks >= p {
				buyLots += q.TotalLots
			}
			return true
		})

		var sellLots int64
		ob.Asks.ForEachLevel0(func(q *OrderQueue) bool {
			if q.PriceTicks <= p {
				sellLots += q.TotalLots
			}
			return true
		})

		exec := minInt64(buyLots, sellLots)
		imbalance := buyLots - sellLots
		if imbalance < 0 {
			imbalance = -imbalance
		}
		distance := int64(0)
		if lastPriceTicks > 0 {
			distance = p - lastPriceTicks
			if distance < 0 {
				distance = -distance
			}
		}

		if !hasBest {
			hasBest = true
			bestPrice = p
			bestExec = exec
			bestImbalance = imbalance
			bestDistance = distance
			continue
		}

		if exec > bestExec {
			bestPrice = p
			bestExec = exec
			bestImbalance = imbalance
			bestDistance = distance
			continue
		}

		if exec == bestExec {
			if imbalance < bestImbalance {
				bestPrice = p
				bestImbalance = imbalance
				bestDistance = distance
				continue
			}
			if imbalance == bestImbalance {
				if lastPriceTicks > 0 && distance < bestDistance {
					bestPrice = p
					bestDistance = distance
					continue
				}
				if lastPriceTicks <= 0 && p < bestPrice {
					bestPrice = p
					continue
				}
			}
		}
	}

	if bestExec == 0 {
		return 0, false
	}
	return bestPrice, true
}

func (ob *OrderBook) MatchAuction(clearingPriceTicks int64) []Trade {
	var trades []Trade

	for {
		bidLevel := ob.Bids.Best()
		askLevel := ob.Asks.Best()
		if bidLevel == nil || askLevel == nil {
			break
		}
		if bidLevel.PriceTicks < clearingPriceTicks || askLevel.PriceTicks > clearingPriceTicks {
			break
		}

		bidOrder := bidLevel.Head()
		askOrder := askLevel.Head()
		if bidOrder == nil || askOrder == nil {
			break
		}

		matchLots := minInt64(bidOrder.RemainingLots, askOrder.RemainingLots)
		ts := ob.now()

		price := ob.priceFromTicks(clearingPriceTicks)
		amount := ob.qtyFromLots(matchLots)

		trade := Trade{
			ID:             ob.nextTradeID(),
			Symbol:         ob.Symbol,
			MakerByID:      askOrder.ID,
			TakerByID:      bidOrder.ID,
			MakerAccountID: askOrder.AccountID,
			TakerAccountID: bidOrder.AccountID,
			Price:          price,
			Amount:         amount,
			Notional:       price.Mul(amount),
			Fee:            decimal.Zero,
			Side:           SideBuy.String(),
			Timestamp:      ts,
			IsBuyerMake:    false,
		}
		trades = append(trades, trade)

		bidOrder.FilledLots += matchLots
		bidOrder.RemainingLots -= matchLots
		bidLevel.TotalLots -= matchLots
		if bidLevel.TotalLots < 0 {
			bidLevel.TotalLots = 0
		}

		askOrder.FilledLots += matchLots
		askOrder.RemainingLots -= matchLots
		askLevel.TotalLots -= matchLots
		if askLevel.TotalLots < 0 {
			askLevel.TotalLots = 0
		}

		if bidOrder.RemainingLots == 0 {
			bidLevel.Pop()
			if bidOrder.Type == TypeIceberg && bidOrder.HiddenRemainingLots > 0 && bidOrder.PeakLots > 0 {
				replenishLots := minInt64(bidOrder.PeakLots, bidOrder.HiddenRemainingLots)
				bidOrder.HiddenRemainingLots -= replenishLots
				bidOrder.RemainingLots = replenishLots
				bidOrder.Status = StatusPartiallyFilled
				bidLevel.Add(bidOrder)
				ob.materializeOrder(bidOrder)
				ob.makerUpdates = append(ob.makerUpdates, bidOrder.Record())
			} else {
				bidOrder.Status = StatusFilled
				ob.materializeOrder(bidOrder)
				ob.makerUpdates = append(ob.makerUpdates, bidOrder.Record())
				delete(ob.Orders, bidOrder.ID)
			}
		} else {
			bidOrder.Status = StatusPartiallyFilled
			ob.materializeOrder(bidOrder)
			ob.makerUpdates = append(ob.makerUpdates, bidOrder.Record())
		}

		if askOrder.RemainingLots == 0 {
			askLevel.Pop()
			if askOrder.Type == TypeIceberg && askOrder.HiddenRemainingLots > 0 && askOrder.PeakLots > 0 {
				replenishLots := minInt64(askOrder.PeakLots, askOrder.HiddenRemainingLots)
				askOrder.HiddenRemainingLots -= replenishLots
				askOrder.RemainingLots = replenishLots
				askOrder.Status = StatusPartiallyFilled
				askLevel.Add(askOrder)
				ob.materializeOrder(askOrder)
				ob.makerUpdates = append(ob.makerUpdates, askOrder.Record())
			} else {
				askOrder.Status = StatusFilled
				ob.materializeOrder(askOrder)
				ob.makerUpdates = append(ob.makerUpdates, askOrder.Record())
				delete(ob.Orders, askOrder.ID)
			}
		} else {
			askOrder.Status = StatusPartiallyFilled
			ob.materializeOrder(askOrder)
			ob.makerUpdates = append(ob.makerUpdates, askOrder.Record())
		}

		if bidLevel.Len() == 0 {
			ob.Bids.RemoveLevel(bidLevel.PriceTicks)
		}
		if askLevel.Len() == 0 {
			ob.Asks.RemoveLevel(askLevel.PriceTicks)
		}
	}

	return trades
}
