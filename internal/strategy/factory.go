package strategy

import (
	"fmt"
)

func NewStrategy(strategyType string, config map[string]interface{}) (Strategy, error) {
	switch strategyType {
	case "ma_cross":
		short, ok1 := config["short_period"].(float64)
		long, ok2 := config["long_period"].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid config for ma_cross: need short_period and long_period")
		}
		return NewMAStrategy(int(short), int(long)), nil
	case "ma_cross_v2":
		short, ok1 := config["short_period"].(float64)
		long, ok2 := config["long_period"].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid config for ma_cross_v2: need short_period and long_period")
		}
		return NewMACrossStrategy(int(short), int(long)), nil
	default:
		return nil, fmt.Errorf("unknown strategy type: %s", strategyType)
	}
}
