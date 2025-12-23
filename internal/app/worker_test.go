package app

import (
	"testing"
)

func TestNormalizeSymbol(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"BTC-USDT", "BTCUSDT"},
		{"btcusdt", "BTCUSDT"},
		{"BTC/USDT", "BTCUSDT"},
		{"ETH_USDT", "ETHUSDT"},
		{"XBT/USD", "XBTUSD"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormalizeSymbol(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeSymbol(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
