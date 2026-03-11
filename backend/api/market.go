package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quant-trader/internal/repo"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListMarketStrategies(c *gin.Context) {
	items, err := h.bizMarketplace.ListPublic(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch marketplace"})
		return
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":          item.ID,
			"price":       item.Price,
			"description": item.Description,
			"metrics":     json.RawMessage(item.Metrics),
			"name":        item.Name,
			"author":      item.AuthorEmail,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) PurchaseStrategy(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	marketItemIDStr := c.Param("id")
	marketItemID, err := strconv.ParseInt(marketItemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.bizMarketplace.Purchase(c.Request.Context(), userID, marketItemID); err != nil {
		if err == repo.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "strategy not found in marketplace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe to strategy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase successful"})
}
