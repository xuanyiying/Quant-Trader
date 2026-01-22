package app

type ingestionTarget struct {
	Exchange string
	Symbol   string
}

func (a *App) ingestionTargets() []ingestionTarget {
	return []ingestionTarget{
		{Exchange: "binance", Symbol: "btcusdt"},
		{Exchange: "okx", Symbol: "BTC-USDT"},
		{Exchange: "bybit", Symbol: "BTCUSDT"},
		{Exchange: "coinbase", Symbol: "BTC-USD"},
		{Exchange: "kraken", Symbol: "XBT/USD"},
	}
}

func (a *App) matchingSymbols() []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, 8)
	for _, t := range a.ingestionTargets() {
		s := NormalizeSymbol(t.Symbol)
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
