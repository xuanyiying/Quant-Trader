package api

import (
	"net/http"
	"strings"
	"time"

	"market-ingestor/internal/engine"
	"market-ingestor/internal/model"
	"market-ingestor/internal/strategy"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewHandler(db *pgxpool.Pool, logger *zap.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
	}
}

// Auth Handlers

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	var userID int64
	err = h.db.QueryRow(c.Request.Context(),
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
		req.Email, string(hash)).Scan(&userID)

	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created", "id": userID})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int64
	var hash string
	err := h.db.QueryRow(c.Request.Context(),
		"SELECT id, password_hash FROM users WHERE email = $1", req.Email).Scan(&userID, &hash)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Data Handlers

func (h *Handler) GetHistoryKLines(c *gin.Context) {
	symbol := strings.ReplaceAll(strings.ToUpper(c.Param("symbol")), "-", "")
	symbol = strings.ReplaceAll(symbol, "/", "")
	period := c.DefaultQuery("period", "1m")

	rows, err := h.db.Query(c.Request.Context(),
		"SELECT symbol, exchange, open, high, low, close, volume, time FROM klines WHERE symbol = $1 AND period = $2 ORDER BY time DESC LIMIT 100",
		symbol, period)
	if err != nil {
		h.logger.Error("failed to query klines", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer rows.Close()

	klines := make([]model.KLine, 0)
	for rows.Next() {
		var k model.KLine
		if err := rows.Scan(&k.Symbol, &k.Exchange, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume, &k.Timestamp); err != nil {
			h.logger.Error("failed to scan kline", zap.Error(err))
			continue
		}
		k.Period = period
		klines = append(klines, k)
	}

	c.JSON(http.StatusOK, klines)
}

func (h *Handler) RunBacktest(c *gin.Context) {
	var req struct {
		Symbol         string                 `json:"symbol" binding:"required"`
		StrategyType   string                 `json:"strategy_type" binding:"required"`
		Config         map[string]interface{} `json:"config"`
		InitialBalance decimal.Decimal        `json:"initial_balance"`
		StartTime      time.Time              `json:"start_time" binding:"required"`
		EndTime        time.Time              `json:"end_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	symbol := strings.ReplaceAll(strings.ToUpper(req.Symbol), "-", "")
	symbol = strings.ReplaceAll(symbol, "/", "")

	// 1. Fetch history data for backtest
	rows, err := h.db.Query(c.Request.Context(),
		"SELECT symbol, exchange, open, high, low, close, volume, time FROM klines WHERE symbol = $1 AND time BETWEEN $2 AND $3 ORDER BY time ASC",
		symbol, req.StartTime, req.EndTime)
	if err != nil {
		h.logger.Error("failed to fetch history for backtest", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch data"})
		return
	}
	defer rows.Close()

	klines := make([]model.KLine, 0)
	for rows.Next() {
		var k model.KLine
		if err := rows.Scan(&k.Symbol, &k.Exchange, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume, &k.Timestamp); err != nil {
			continue
		}
		klines = append(klines, k)
	}

	// 2. Setup Strategy
	strat, err := strategy.NewStrategy(req.StrategyType, req.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Run Backtest
	tester := engine.NewBacktester(strat, req.InitialBalance)
	report := tester.Run(klines)

	// Optional: Save backtest run to DB
	// ... (omitted for brevity)

	c.JSON(http.StatusOK, report)
}
