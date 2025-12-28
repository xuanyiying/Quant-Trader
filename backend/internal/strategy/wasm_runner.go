package strategy

import (
	"context"
	"fmt"
	"quant-trader/internal/model"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type WasmRunner struct {
	runtime wazero.Runtime
	mod     api.Module
}

func NewWasmRunner(ctx context.Context, wasmCode []byte) (*WasmRunner, error) {
	r := wazero.NewRuntime(ctx)

	// Compile and instantiate the module
	mod, err := r.Instantiate(ctx, wasmCode)
	if err != nil {
		r.Close(ctx)
		return nil, fmt.Errorf("failed to instantiate wasm: %w", err)
	}

	return &WasmRunner{
		runtime: r,
		mod:     mod,
	}, nil
}

func (r *WasmRunner) OnCandle(candle model.KLine) Action {
	// 1. Get the handle to the OnCandle function in Wasm
	onCandle := r.mod.ExportedFunction("OnCandle")
	if onCandle == nil {
		return ActionHold
	}

	// 2. Map Go candle to Wasm memory (simplified example)
	price, _ := candle.Close.Float64()

	// 3. Call the function
	results, err := onCandle.Call(context.Background(), api.EncodeF64(price))
	if err != nil {
		return ActionHold
	}

	// 4. Map the return value back to Action
	if len(results) > 0 {
		// Define a mapping between Wasm integer returns and our string Actions
		switch results[0] {
		case 1:
			return ActionBuy
		case 2:
			return ActionSell
		default:
			return ActionHold
		}
	}

	return ActionHold
}

func (r *WasmRunner) OnTrade(trade model.Trade) Action {
	// 1. Get the handle to the OnTrade function in Wasm
	onTrade := r.mod.ExportedFunction("OnTrade")
	if onTrade == nil {
		// If not implemented in Wasm, just return Hold
		return ActionHold
	}

	// 2. Map Go trade to Wasm memory
	price, _ := trade.Price.Float64()

	// 3. Call the function (using Background context since the interface doesn't pass one)
	results, err := onTrade.Call(context.Background(), api.EncodeF64(price))
	if err != nil {
		return ActionHold
	}

	// 4. Map the return value back to Action
	if len(results) > 0 {
		switch results[0] {
		case 1:
			return ActionBuy
		case 2:
			return ActionSell
		default:
			return ActionHold
		}
	}

	return ActionHold
}

func (r *WasmRunner) Close(ctx context.Context) error {
	return r.runtime.Close(ctx)
}
