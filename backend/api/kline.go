package api

import (
	"net/http"
	"strings"
	"time"

	"quant-trader/internal/biz"
	"quant-trader/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (h *Handler) GetHistoryKLines(c *gin.Context) {
	symbol := strings.ReplaceAll(strings.ToUpper(c.Param("symbol")), "-", "")
	symbol = strings.ReplaceAll(symbol, "/", "")
	period := c.DefaultQuery("period", model.Period1m)

	klines, err := h.bizKline.GetLatest(c.Request.Context(), symbol, period, model.DefaultHistoryLimit)
	if err != nil {
		h.logger.Error("failed to query klines", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, klines)
}

func (h *Handler) RunBacktest(c *gin.Context) {
	var req struct {
		Symbol         string                 `json:"symbol" binding:"required"`
		StrategyType   string                 `json:"strategy_type" binding:"required"`
		Config         map[string]interface{} `json:"config"`
		InitialBalance decimal.Decimal        `json:"initial_balance"`
		StartTime      time.Time              `json:"start_time" binding:"required"`
		EndTime        time.Time              `json:"end_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	symbol := strings.ReplaceAll(strings.ToUpper(req.Symbol), "-", "")
	symbol = strings.ReplaceAll(symbol, "/", "")

	report, err := h.bizKline.RunBacktest(c.Request.Context(), biz.BacktestInput{
		Symbol:         symbol,
		StrategyType:   req.StrategyType,
		Config:         req.Config,
		InitialBalance: req.InitialBalance,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
	})

	if err != nil {
		h.logger.Error("failed to run backtest", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *Handler) TriggerBackfill(c *gin.Context) {
	var req struct {
		Exchange  string    `json:"exchange" binding:"required"`
		Symbol    string    `json:"symbol" binding:"required"`
		StartTime time.Time `json:"start_time" binding:"required"`
		EndTime   time.Time `json:"end_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 Biz 层启动异步补全任务
	h.bizBackfill.BackfillExchange(c.Request.Context(), strings.ToLower(req.Exchange), strings.ToUpper(req.Symbol), req.StartTime, req.EndTime)

	c.JSON(http.StatusAccepted, gin.H{"message": "backfill task started"})
}
