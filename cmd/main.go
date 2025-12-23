package main

import (
	"context"
	"log"

	"market-ingestor/internal/app"
)

func main() {
	// Create application instance
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}

	// Initialize application (DB, NATS, etc.)
	ctx := context.Background()
	if err := application.Init(ctx); err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	// Run application
	if err := application.Run(ctx); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
