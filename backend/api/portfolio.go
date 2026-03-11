package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPortfolios(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	portfolios, err := h.bizPortfolio.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch portfolios"})
		return
	}

	result := make([]map[string]interface{}, len(portfolios))
	for i, p := range portfolios {
		result[i] = map[string]interface{}{
			"id":         p.ID,
			"name":       p.Name,
			"created_at": p.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreatePortfolio(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio, err := h.bizPortfolio.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create portfolio"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": portfolio.ID})
}

func (h *Handler) GetPortfolioReport(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	report, err := h.bizPortfolio.GetReport(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}
