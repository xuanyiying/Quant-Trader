package api

import (
	"fmt"
	"net/http"
	"time"

	"quant-trader/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) ListAPIKeys(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	keys, err := h.apiKey.GetByUserID(c.Request.Context(), userID)
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

	keyID := fmt.Sprintf("qt_%d_%d", userID, time.Now().Unix())
	rawSecret := fmt.Sprintf("sec_%d_%d", userID, time.Now().UnixNano())
	hash, _ := bcrypt.GenerateFromPassword([]byte(rawSecret), bcrypt.DefaultCost)

	key := &model.APIKey{
		UserID:    userID,
		Name:      req.Name,
		KeyPrefix: keyID,
		KeyHash:   string(hash),
		IsActive:  true,
	}

	err := h.apiKey.Create(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create API key"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"key_id":     keyID,
		"key_secret": rawSecret,
		"message":    "Store the secret safely, it will not be shown again",
	})
}
