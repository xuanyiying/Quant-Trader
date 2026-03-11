package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListAPIKeys(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	keys, err := h.bizAPIKey.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch API keys"})
		return
	}

	result := make([]map[string]interface{}, len(keys))
	for i, k := range keys {
		result[i] = map[string]interface{}{
			"id":         k.ID,
			"key_id":     k.KeyPrefix,
			"name":       k.Name,
			"is_active":  k.IsActive,
			"created_at": k.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateAPIKey(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.bizAPIKey.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create API key"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"key_id":     result.KeyID,
		"key_secret": result.KeySecret,
		"message":    "Store the secret safely, it will not be shown again",
	})
}
