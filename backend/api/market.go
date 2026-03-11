package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (h *Handler) ListMarketStrategies(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(),
		`SELECT m.id, m.price, m.description, m.performance_metrics, s.name, u.email as author
		 FROM strategy_market m
		 JOIN strategies s ON m.strategy_id = s.id
		 JOIN users u ON m.owner_id = u.id
		 WHERE m.is_public = TRUE`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch marketplace"})
		return
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var (
			id      int64
			price   decimal.Decimal
			desc    string
			metrics []byte
			sname   string
			author  string
		)
		if err := rows.Scan(&id, &price, &desc, &metrics, &sname, &author); err != nil {
			continue
		}
		items = append(items, map[string]interface{}{
			"id":          id,
			"price":       price,
			"description": desc,
			"metrics":     json.RawMessage(metrics),
			"name":        sname,
			"author":      author,
		})
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) PurchaseStrategy(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	marketItemID := c.Param("id")

	// 1. Check if the strategy exists and get its price
	var price decimal.Decimal
	err := h.db.QueryRow(c.Request.Context(),
		"SELECT price FROM strategy_market WHERE id = $1 AND is_public = TRUE",
		marketItemID).Scan(&price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "strategy not found in marketplace"})
		return
	}

	// 2. If price > 0, verify payment (In a real scenario, integrate with Stripe or check user balance)
	if price.GreaterThan(decimal.Zero) {
		// For now, we still just grant access but we've added the price check.
		// In production, you'd verify a payment_intent_id or deduct from internal credits.
		h.logger.Info("purchasing paid strategy", zap.Int64("user_id", userID), zap.String("item_id", marketItemID), zap.String("price", price.String()))
	}

	// 3. Grant access
	_, err = h.db.Exec(c.Request.Context(),
		"INSERT INTO strategy_purchases (user_id, market_item_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID, marketItemID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe to strategy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "strategy subscribed successfully"})
}
