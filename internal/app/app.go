package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"market-ingestor/api"
	"market-ingestor/internal/config"
	"market-ingestor/internal/infrastructure"
	"market-ingestor/internal/processor"
	"market-ingestor/internal/push"
	"market-ingestor/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// App defines the application structure and its dependencies
type App struct {
	Config      *config.Config
	Logger      *zap.Logger
	DB          *pgxpool.Pool
	NC          *nats.Conn
	JS          nats.JetStreamContext
	PushGateway *push.PushGateway
	HTTPServer  *http.Server
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	infrastructure.Init()
	logger := infrastructure.Logger

	return &App{
		Config: &cfg,
		Logger: logger,
	}, nil
}

// Init initializes all application components
func (a *App) Init(ctx context.Context) error {
	// 1. Database
	dbPool, err := pgxpool.Connect(ctx, a.Config.DB_DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.DB = dbPool

	if err := a.initDatabase(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// 2. NATS
	nc, js, err := infrastructure.InitNATS(a.Config.NatsURL, a.Logger)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	a.NC = nc
	a.JS = js

	// 3. Services
	a.PushGateway = push.NewPushGateway(js, a.Logger)

	return nil
}

// Run starts the application services and the HTTP server
func (a *App) Run(ctx context.Context) error {
	// Start Persistence Service
	tradeSaver := storage.NewBatchSaver(a.DB, a.Logger, 1*time.Second, 1000)
	klineSaver := storage.NewKlineSaver(a.DB, a.Logger, 1*time.Second, 100)
	a.startPersistenceService(tradeSaver, klineSaver)

	// Start Stream Processor
	klineProcessor := processor.NewKlineProcessor(a.JS, a.Logger)
	if err := klineProcessor.Run(ctx); err != nil {
		return fmt.Errorf("failed to start kline processor: %w", err)
	}

	// Start Ingestion Worker
	a.startIngestionWorker(ctx)

	// Setup HTTP Server
	a.HTTPServer = &http.Server{
		Addr:    ":" + a.Config.Port,
		Handler: a.setupRouter(),
	}

	go func() {
		a.Logger.Info("starting http server", zap.String("port", a.Config.Port))
		if err := a.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	return a.waitForShutdown()
}

// waitForShutdown handles graceful shutdown signals
func (a *App) waitForShutdown() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	a.Logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	a.NC.Close()
	a.DB.Close()

	return nil
}

// initDatabase runs the database initialization script
func (a *App) initDatabase(ctx context.Context) error {
	sqlFile := "scripts/init.sql"
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read init script: %w", err)
	}

	_, err = a.DB.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("failed to execute init script: %w", err)
	}

	a.Logger.Info("database initialized successfully")
	return nil
}

// setupRouter configures the Gin router and its routes
func (a *App) setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	apiHandler := api.NewHandler(a.DB, a.Logger)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", apiHandler.Register)
		v1.POST("/login", apiHandler.Login)
		v1.GET("/klines/:symbol", apiHandler.GetHistoryKLines)
	}

	r.StaticFile("/", "./static/index.html")
	r.Static("/static", "./static")

	protected := r.Group("/api/v1")
	protected.Use(api.AuthMiddleware())
	{
		protected.POST("/backtest", apiHandler.RunBacktest)
	}

	r.GET("/ws", func(c *gin.Context) {
		a.PushGateway.ServeHTTP(c.Writer, c.Request)
	})

	return r
}
