package biz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test error definitions
func TestErrorDefinitions(t *testing.T) {
	assert.NotNil(t, ErrInvalidOrderSide)
	assert.NotNil(t, ErrInvalidOrderType)
	assert.Equal(t, "invalid order side", ErrInvalidOrderSide.Error())
	assert.Equal(t, "invalid order type", ErrInvalidOrderType.Error())
}

// Test CreateOrderInput validation
func TestCreateOrderInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateOrderInput
		wantErr error
	}{
		{
			name: "invalid side - empty",
			input: CreateOrderInput{
				Symbol: "BTCUSDT",
				Side:   "",
				Type:   "market",
			},
			wantErr: ErrInvalidOrderSide,
		},
		{
			name: "invalid side - unknown",
			input: CreateOrderInput{
				Symbol: "BTCUSDT",
				Side:   "unknown",
				Type:   "market",
			},
			wantErr: ErrInvalidOrderSide,
		},
		{
			name: "invalid type - empty",
			input: CreateOrderInput{
				Symbol: "BTCUSDT",
				Side:   "buy",
				Type:   "",
			},
			wantErr: ErrInvalidOrderType,
		},
		{
			name: "invalid type - unknown",
			input: CreateOrderInput{
				Symbol: "BTCUSDT",
				Side:   "buy",
				Type:   "unknown",
			},
			wantErr: ErrInvalidOrderType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOrderInput(tt.input)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// validateOrderInput 验证订单输入
func validateOrderInput(input CreateOrderInput) error {
	if input.Side != "buy" && input.Side != "sell" {
		return ErrInvalidOrderSide
	}
	if input.Type != "market" && input.Type != "limit" {
		return ErrInvalidOrderType
	}
	return nil
}