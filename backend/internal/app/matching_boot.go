package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"quant-trader/internal/matching"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (a *App) initMatching(ctx context.Context) error {
	if a.Config == nil || !a.Config.MatchingEnabled {
		return nil
	}

	ttl, err := time.ParseDuration(a.Config.MatchingLeaseTTL)
	if err != nil {
		return fmt.Errorf("invalid MATCHING_LEASE_TTL: %w", err)
	}
	refresh, err := time.ParseDuration(a.Config.MatchingLeaseRefresh)
	if err != nil {
		return fmt.Errorf("invalid MATCHING_LEASE_REFRESH: %w", err)
	}
	snapshotEvery, err := time.ParseDuration(a.Config.MatchingSnapshotEvery)
	if err != nil {
		return fmt.Errorf("invalid MATCHING_SNAPSHOT_EVERY: %w", err)
	}

	owner := a.Config.MatchingOwner
	if owner == "" {
		host, _ := os.Hostname()
		owner = fmt.Sprintf("%s-%d", host, os.Getpid())
	}

	engine := matching.NewEngine()
	wal, err := matching.NewJetStreamWAL(a.JS)
	if err != nil {
		return fmt.Errorf("failed to init matching WAL: %w", err)
	}
	snapshotStore, err := matching.NewJetStreamSnapshotStore(a.JS)
	if err != nil {
		return fmt.Errorf("failed to init matching snapshot store: %w", err)
	}
	lease, err := matching.NewJetStreamLease(a.JS, ttl)
	if err != nil {
		return fmt.Errorf("failed to init matching lease: %w", err)
	}

	a.MatchEngine = engine
	a.MatchWAL = wal
	a.MatchSnapshot = snapshotStore
	a.MatchLease = lease

	matching.StartWALPump(ctx, a.MatchEngine.Events(), a.MatchWAL)

	for _, symbol := range a.matchingSymbols() {
		a.startSymbolLeader(ctx, symbol, owner, refresh, snapshotEvery)
	}

	return nil
}

func (a *App) startSymbolLeader(ctx context.Context, symbol string, owner string, refresh time.Duration, snapshotEvery time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			release, err := a.MatchLease.Acquire(ctx, symbol, owner, refresh)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			symCtx, cancel := context.WithCancel(ctx)
			if err := a.restoreSymbol(symCtx, symbol); err != nil {
				cancel()
				_ = release()
				a.Logger.Error("failed to restore matching symbol", zap.String("symbol", symbol), zap.Error(err))
				time.Sleep(2 * time.Second)
				continue
			}

			matching.StartSnapshotPump(symCtx, a.MatchEngine, a.MatchSnapshot, symbol, snapshotEvery)
			a.Logger.Info("matching symbol leader active", zap.String("symbol", symbol), zap.String("owner", owner))

			<-symCtx.Done()
			cancel()
			_ = release()
		}
	}()
}

func (a *App) restoreSymbol(ctx context.Context, symbol string) error {
	a.MatchEngine.AddOrderBook(symbol)
	_ = a.MatchEngine.SetRiskManager(symbol, matching.NewAbuseRiskManager(matching.AbuseConfig{}))

	snap, trades, err := matching.RecoverSymbol(a.MatchSnapshot, a.MatchWAL, symbol)
	if err != nil {
		return err
	}
	if len(trades) > 0 {
		last := trades[len(trades)-1]
		snap.LastPrice = last.Price
		if snap.Config.TickSize.GreaterThan(decimal.Zero) {
			q := last.Price.Div(snap.Config.TickSize).Round(0)
			snap.LastPriceTicks = q.IntPart()
		}
	}

	if err := a.MatchEngine.Restore(symbol, snap); err != nil {
		return err
	}
	return nil
}
