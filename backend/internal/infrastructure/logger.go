package infrastructure

import (
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func Init() {
	Logger, _ = zap.NewProduction()
	Logger.Info("infrastructure initialized")
}
