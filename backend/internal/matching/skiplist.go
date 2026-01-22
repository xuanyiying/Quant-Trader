package matching

import (
	"math/rand"
)

const (
	MaxLevel    = 32
	Probability = 0.5
)

// OrderQueue is a FIFO queue of orders at a specific price level
type OrderQueue struct {
	PriceTicks int64
	TotalLots  int64
	head       *orderNode
	tail       *orderNode
	index      map[string]*orderNode
	length     int
}

type orderNode struct {
	order *Order
	prev  *orderNode
	next  *orderNode
}

func NewOrderQueue(priceTicks int64) *OrderQueue {
	return &OrderQueue{
		PriceTicks: priceTicks,
		index:      make(map[string]*orderNode),
	}
}

func (q *OrderQueue) Add(order *Order) {
	node := &orderNode{order: order}
	if q.tail == nil {
		q.head = node
		q.tail = node
	} else {
		node.prev = q.tail
		q.tail.next = node
		q.tail = node
	}
	q.index[order.ID] = node
	q.length++
	q.TotalLots += order.RemainingLots
}

func (q *OrderQueue) Remove(orderID string) *Order {
	node, ok := q.index[orderID]
	if !ok {
		return nil
	}
	delete(q.index, orderID)
	q.length--
	q.TotalLots -= node.order.RemainingLots

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		q.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		q.tail = node.prev
	}
	node.prev = nil
	node.next = nil
	return node.order
}

func (q *OrderQueue) Head() *Order {
	if q.head == nil {
		return nil
	}
	return q.head.order
}

func (q *OrderQueue) Pop() *Order {
	if q.head == nil {
		return nil
	}
	node := q.head
	q.head = node.next
	if q.head != nil {
		q.head.prev = nil
	} else {
		q.tail = nil
	}
	delete(q.index, node.order.ID)
	q.length--
	q.TotalLots -= node.order.RemainingLots
	node.next = nil
	node.prev = nil
	return node.order
}

func (q *OrderQueue) Len() int {
	return q.length
}

// SkipListNode represents a node in the SkipList
type SkipListNode struct {
	Value *OrderQueue
	Next  []*SkipListNode
}

// SkipList implementation for Price Levels
type SkipList struct {
	Head       *SkipListNode
	Level      int
	Descending bool // true for Bids (High->Low), false for Asks (Low->High)
}

func NewSkipList(descending bool) *SkipList {
	return &SkipList{
		Head: &SkipListNode{
			Next: make([]*SkipListNode, MaxLevel),
		},
		Level:      1,
		Descending: descending,
	}
}

func (sl *SkipList) randomLevel() int {
	level := 1
	for rand.Float64() < Probability && level < MaxLevel {
		level++
	}
	return level
}

// less returns true if p1 should come before p2
func (sl *SkipList) less(p1, p2 int64) bool {
	if sl.Descending {
		return p1 > p2
	}
	return p1 < p2
}

func (sl *SkipList) Insert(priceTicks int64, order *Order) *OrderQueue {
	update := make([]*SkipListNode, MaxLevel)
	current := sl.Head

	for i := sl.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil {
			nextPrice := current.Next[i].Value.PriceTicks
			if nextPrice == priceTicks {
				// Found existing level
				current.Next[i].Value.Add(order)
				return current.Next[i].Value
			}
			if sl.less(nextPrice, priceTicks) {
				current = current.Next[i]
			} else {
				break
			}
		}
		update[i] = current
	}

	// If we are here, we need to create a new level/node
	// Double check level 0 just in case
	if current.Next[0] != nil && current.Next[0].Value.PriceTicks == priceTicks {
		current.Next[0].Value.Add(order)
		return current.Next[0].Value
	}

	newLevel := sl.randomLevel()
	if newLevel > sl.Level {
		for i := sl.Level; i < newLevel; i++ {
			update[i] = sl.Head
		}
		sl.Level = newLevel
	}

	newQueue := NewOrderQueue(priceTicks)
	newQueue.Add(order)
	newNode := &SkipListNode{
		Value: newQueue,
		Next:  make([]*SkipListNode, newLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.Next[i] = update[i].Next[i]
		update[i].Next[i] = newNode
	}

	return newQueue
}

// Get returns the OrderQueue for a specific price
func (sl *SkipList) Get(priceTicks int64) *OrderQueue {
	current := sl.Head
	for i := sl.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil {
			nextPrice := current.Next[i].Value.PriceTicks
			if nextPrice == priceTicks {
				return current.Next[i].Value
			}
			if sl.less(nextPrice, priceTicks) {
				current = current.Next[i]
			} else {
				break
			}
		}
	}
	return nil
}

// RemoveLevel removes the entire level at the given price
func (sl *SkipList) RemoveLevel(priceTicks int64) bool {
	update := make([]*SkipListNode, MaxLevel)
	current := sl.Head

	for i := sl.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil {
			nextPrice := current.Next[i].Value.PriceTicks
			if sl.less(nextPrice, priceTicks) {
				current = current.Next[i]
			} else {
				break
			}
		}
		update[i] = current
	}

	target := current.Next[0]
	if target != nil && target.Value.PriceTicks == priceTicks {
		for i := 0; i < sl.Level; i++ {
			if update[i].Next[i] != target {
				break
			}
			update[i].Next[i] = target.Next[i]
		}
		for sl.Level > 1 && sl.Head.Next[sl.Level-1] == nil {
			sl.Level--
		}
		return true
	}
	return false
}

// Remove removes an order from the queue at the given price.
// If the queue becomes empty, the node is removed from the SkipList.
func (sl *SkipList) Remove(priceTicks int64, orderID string) bool {
	update := make([]*SkipListNode, MaxLevel)
	current := sl.Head

	for i := sl.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil {
			nextPrice := current.Next[i].Value.PriceTicks
			if sl.less(nextPrice, priceTicks) {
				current = current.Next[i]
			} else {
				break
			}
		}
		update[i] = current
	}

	target := current.Next[0]
	if target != nil && target.Value.PriceTicks == priceTicks {
		// Found the price level
		removedOrder := target.Value.Remove(orderID)
		if removedOrder != nil {
			// If queue is empty, remove the node
			if target.Value.Len() == 0 {
				for i := 0; i < sl.Level; i++ {
					if update[i].Next[i] != target {
						break
					}
					update[i].Next[i] = target.Next[i]
				}
				// Adjust level
				for sl.Level > 1 && sl.Head.Next[sl.Level-1] == nil {
					sl.Level--
				}
			}
			return true
		}
	}
	return false
}

// Best returns the top OrderQueue (Best Bid or Best Ask)
func (sl *SkipList) Best() *OrderQueue {
	if sl.Head.Next[0] != nil {
		return sl.Head.Next[0].Value
	}
	return nil
}

func (sl *SkipList) ForEachLevel0(fn func(q *OrderQueue) bool) {
	for node := sl.Head.Next[0]; node != nil; node = node.Next[0] {
		if !fn(node.Value) {
			return
		}
	}
}
