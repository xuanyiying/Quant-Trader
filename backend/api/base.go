package api

import (
	"os"

	"quant-trader/internal/analytics"
	"quant-trader/internal/biz"
	"quant-trader/internal/payment"
	"quant-trader/internal/repo"
	"quant-trader/internal/risk"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	db     *pgxpool.Pool
	logger *zap.Logger
	stripe *payment.StripeService

	bizAuth         *biz.Auth
	bizPaper        *biz.PaperTrading
	bizPortfolio    *biz.Portfolio
	bizSubscription *biz.Subscription
	bizAPIKey       *biz.APIKey
	bizAlert        *biz.Alert
	bizKline        *biz.Kline
	bizMarketplace  *biz.Marketplace
}

func NewHandler(db *pgxpool.Pool, gormDB *gorm.DB, logger *zap.Logger) *Handler {
	stripeKey := os.Getenv("STRIPE_API_KEY")

	h := &Handler{
		db:     db,
		logger: logger,
		stripe: payment.NewStripeService(db, logger, stripeKey),
	}

	repoUser := repo.NewUser(gormDB)
	repoAPIKey := repo.NewAPIKey(gormDB)
	repoPaper := repo.NewPaperAccount(gormDB)
	repoPosition := repo.NewPaperPosition(gormDB)
	repoOrder := repo.NewPaperOrder(gormDB)
	repoAlert := repo.NewAlert(gormDB)
	repoPortfolio := repo.NewPortfolio(gormDB)
	repoSubscription := repo.NewSubscription(gormDB)
	repoKline := repo.NewKline(gormDB)
	repoMarketplace := repo.NewMarketplace(gormDB)

	riskMgr := risk.NewRiskManager(db, logger)
	analyticsSvc := analytics.NewAnalyticsService(db, logger)

	h.bizAuth = biz.NewAuth(repoUser, logger)
	h.bizPaper = biz.NewPaperTrading(repoPaper, repoPosition, repoOrder, riskMgr, logger)
	h.bizPortfolio = biz.NewPortfolio(repoPortfolio, analyticsSvc, logger)
	h.bizSubscription = biz.NewSubscription(repoSubscription, logger)
	h.bizAPIKey = biz.NewAPIKey(repoAPIKey, logger)
	h.bizAlert = biz.NewAlert(repoAlert, logger)
	h.bizKline = biz.NewKline(repoKline, logger)
	h.bizMarketplace = biz.NewMarketplace(repoMarketplace, logger)

	return h
}
