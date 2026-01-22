package matching

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Lease interface {
	Acquire(ctx context.Context, symbol string, owner string, refresh time.Duration) (func() error, error)
}

type MemoryLease struct {
	mu    sync.Mutex
	owner map[string]string
}

func NewMemoryLease() *MemoryLease {
	return &MemoryLease{owner: make(map[string]string)}
}

func (l *MemoryLease) Acquire(ctx context.Context, symbol string, owner string, refresh time.Duration) (func() error, error) {
	l.mu.Lock()
	if cur, ok := l.owner[symbol]; ok && cur != "" {
		l.mu.Unlock()
		return nil, fmt.Errorf("lease already held")
	}
	l.owner[symbol] = owner
	l.mu.Unlock()

	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(refresh)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
			case <-stop:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return func() error {
		close(stop)
		l.mu.Lock()
		defer l.mu.Unlock()
		if cur, ok := l.owner[symbol]; ok && cur == owner {
			delete(l.owner, symbol)
		}
		return nil
	}, nil
}
