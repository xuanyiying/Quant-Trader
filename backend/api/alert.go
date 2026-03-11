package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"quant-trader/internal/biz"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (h *Handler) GetAlerts(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	alerts, err := h.bizAlert.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch alerts"})
		return
	}

	result := make([]map[string]interface{}, len(alerts))
	for i, a := range alerts {
		result[i] = map[string]interface{}{
			"id":             a.ID,
			"symbol":         a.Symbol,
			"condition_type": a.Condition,
			"target_value":   a.Price,
			"is_triggered":   a.IsTriggered,
			"created_at":     a.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateAlert(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var req struct {
		Symbol        string          `json:"symbol" binding:"required"`
		ConditionType string          `json:"condition_type" binding:"required"`
		TargetValue   decimal.Decimal `json:"target_value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.bizAlert.Create(c.Request.Context(), userID, strings.ToUpper(req.Symbol), req.ConditionType, req.TargetValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create alert"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": alert.ID})
}

func (h *Handler) DeleteAlert(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	alertID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert id"})
		return
	}

	if err := h.bizAlert.Delete(c.Request.Context(), userID, alertID); err != nil {
		if errors.Is(err, biz.ErrAlertNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert deleted"})
}
