package api

import (
	"net/http"
	"strings"

	"quant-trader/internal/biz"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (h *Handler) GetPaperAccount(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	account, err := h.bizPaper.GetOrCreateAccount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance})
}

func (h *Handler) ResetPaperAccount(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	if err := h.bizPaper.ResetAccount(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": decimal.NewFromFloat(100000),
		"message": "account reset successful",
	})
}

func (h *Handler) CreatePaperOrder(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var req struct {
		Symbol string          `json:"symbol" binding:"required"`
		Side   string          `json:"side" binding:"required"`
		Type   string          `json:"type" binding:"required"`
		Price  decimal.Decimal `json:"price"`
		Qty    decimal.Decimal `json:"qty" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.bizPaper.CreateOrder(c.Request.Context(), userID, biz.CreateOrderInput{
		Symbol: strings.ToUpper(req.Symbol),
		Side:   req.Side,
		Type:   req.Type,
		Price:  req.Price,
		Qty:    req.Qty,
	})

	if err != nil {
		switch err {
		case biz.ErrInvalidOrderSide, biz.ErrInvalidOrderType:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": order.ID})
}

func (h *Handler) GetPaperPositions(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	positions, err := h.bizPaper.GetPositions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch positions"})
		return
	}

	result := make([]map[string]interface{}, len(positions))
	for i, p := range positions {
		result[i] = map[string]interface{}{
			"symbol":    p.Symbol,
			"qty":       p.Qty,
			"avg_price": p.AvgPrice,
		}
	}

	c.JSON(http.StatusOK, result)
}
