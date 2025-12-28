package middleware

import (
	"net/http"

	"quant-trader/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SubscriptionMiddleware(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		var tierName string
		var maxSymbols int
		err := db.QueryRow(c.Request.Context(),
			`SELECT t.name, t.max_symbols FROM user_subscriptions s 
			 JOIN subscription_tiers t ON s.tier_id = t.id 
			 WHERE s.user_id = $1 AND s.status = $2`, userID, model.SubscriptionStatusActive).Scan(&tierName, &maxSymbols)

		if err != nil {
			// Default to Free tier if no record found
			tierName = model.TierFree
			maxSymbols = 1
		}

		c.Set("tier", tierName)
		c.Set("maxSymbols", maxSymbols)
		c.Next()
	}
}
