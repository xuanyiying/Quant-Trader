package api

import (
	"net/http"
	"strconv"
	"strings"

	"quant-trader/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (h *Handler) GetAlerts(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	alerts, err := h.alert.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch alerts"})
		return
	}

	result := make([]map[string]interface{}, len(alerts))
	for i, a := range alerts {
		result[i] = map[string]interface{}{
			"id":              a.ID,
			"symbol":          a.Symbol,
			"condition_type":  a.Condition,
			"target_value":    a.Price,
			"is_triggered":    a.IsTriggered,
			"created_at":      a.CreatedAt,
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

	alert := &model.Alert{
		UserID:     userID,
		Symbol:     strings.ToUpper(req.Symbol),
		Condition:  req.ConditionType,
		Price:      req.TargetValue,
		IsTriggered: false,
	}

	err := h.alert.Create(c.Request.Context(), alert)
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

	alert, err := h.alert.GetByID(c.Request.Context(), alertID)
	if err != nil || alert.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
		return
	}

	err = h.alert.Delete(c.Request.Context(), alertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert deleted"})
}
