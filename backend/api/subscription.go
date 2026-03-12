package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"quant-trader/internal/biz"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	info, err := h.bizSubscription.GetSubscription(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, biz.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tier_name":   info.TierName,
		"max_symbols": info.MaxSymbols,
		"status":      info.Status,
		"expires_at":  info.ExpiresAt,
	})
}

func (h *Handler) CreateCheckoutSession(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var req struct {
		PriceID string `json:"price_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.stripe.CreateCheckoutSession(userID, req.PriceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create checkout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handler) HandleStripeWebhook(c *gin.Context) {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	err = h.stripe.HandleWebhook(payload, sigHeader, endpointSecret)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
