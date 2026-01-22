I will implement a high-performance Matching Engine in `backend/internal/matching` following your architectural specifications.

### **Plan: High-Performance Matching Engine Implementation**

**1. Data Structures & Core Models**
- Create `backend/internal/matching/model.go` to define:
  - `Order`: Fixed-size structure (where possible) with ID, Price, Size, Side, Type (Limit/Market/IOC/FOK/Iceberg), Timestamp.
  - `Trade`: Execution report.
  - `Side` & `OrderType` enums.
  - **Optimization**: Use `decimal.Decimal` (from existing dep) or `int64` (micros/nanos) for price/size to ensure precision and speed. I will use `decimal.Decimal` to match existing codebase conventions.

**2. Efficient OrderBook Container (SkipList)**
- Create `backend/internal/matching/skiplist.go`.
- Implement a **Skip List** to maintain Price Levels.
  - **Why**: You requested RB-Tree or SkipList. SkipLists are easier to implement concurrency-safe (if needed later) and offer O(log N) insertion/deletion, which is critical for 100k+ TPS compared to O(N) sorted slices.
  - **Structure**: Nodes represent a `PriceLevel`. Each `PriceLevel` contains a FIFO queue of Orders (for Time priority).

**3. OrderBook Logic**
- Create `backend/internal/matching/orderbook.go`.
- **Dual-Side Management**: Maintain `bids` (Max-Heap/Desc SkipList) and `asks` (Min-Heap/Asc SkipList).
- **Matching Loop**:
  - **Limit Order**: Check opposite book. If cross (Bid >= BestAsk or Ask <= BestBid), match. Repeat until filled or no match. Add remainder to book.
  - **Market Order**: Eat liquidity from best price until filled.
  - **Iceberg**: Only display "visible" amount. When filled, reload from "hidden" and re-queue at new time priority (or same, depending on policy - usually new time priority to be fair).
  - **FOK/IOC**: Special handling to cancel immediately if not filled.

**4. Matching Engine Core**
- Create `backend/internal/matching/engine.go`.
- **Sharding**: `map[string]*OrderBook` to manage multiple symbols.
- **Concurrency**: A "Processor" pattern where each symbol (or group of symbols) runs on a dedicated goroutine (Single Writer Principle) to avoid locks within the critical matching path, achieving sub-millisecond latency.
- **Input/Output**: Channels for `SubmitOrder`, `CancelOrder` -> Output Channel for `Events` (Trades, OrderUpdates).

**5. Verification & Testing**
- Create `backend/internal/matching/engine_test.go`.
- **Unit Test**: Verify Price-Time priority (e.g., earlier order gets filled first).
- **Scenario Test**:
  - Place Sell 100 @ 10.0.
  - Place Buy 50 @ 10.0 -> Match 50.
  - Place Buy 60 @ 10.0 -> Match remaining 50, leave 10 on book.
- **Benchmark**: Simple TPS test to verify performance goals.

### **Next Steps**
1.  Create the `backend/internal/matching` directory.
2.  Implement the SkipList and Model definitions.
3.  Implement the OrderBook matching logic.
4.  Implement the Engine wrapper and Tests.
