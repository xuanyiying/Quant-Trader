package api

import (
	"os"

	"quant-trader/internal/analytics"
	"quant-trader/internal/payment"
	"quant-trader/internal/repo"
	"quant-trader/internal/risk"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	db        *pgxpool.Pool
	gormDB    *gorm.DB
	logger    *zap.Logger
	risk      *risk.RiskManager
	analytics *analytics.AnalyticsService
	stripe    *payment.StripeService

	user         *repo.User
	apiKey       *repo.APIKey
	paper        *repo.PaperAccount
	position     *repo.PaperPosition
	alert        *repo.Alert
	portfolio    *repo.Portfolio
	subscription *repo.Subscription
}

func NewHandler(db *pgxpool.Pool, gormDB *gorm.DB, logger *zap.Logger) *Handler {
	stripeKey := os.Getenv("STRIPE_API_KEY")

	h := &Handler{
		db:        db,
		gormDB:    gormDB,
		logger:    logger,
		risk:      risk.NewRiskManager(db, logger),
		analytics: analytics.NewAnalyticsService(db, logger),
		stripe:    payment.NewStripeService(db, logger, stripeKey),
	}

	h.user = repo.NewUser(gormDB)
	h.apiKey = repo.NewAPIKey(gormDB)
	h.paper = repo.NewPaperAccount(gormDB)
	h.position = repo.NewPaperPosition(gormDB)
	h.alert = repo.NewAlert(gormDB)
	h.portfolio = repo.NewPortfolio(gormDB)
	h.subscription = repo.NewSubscription(gormDB)

	return h
}
