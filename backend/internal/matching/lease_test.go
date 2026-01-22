package matching

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryLease_Takeover(t *testing.T) {
	lease := NewMemoryLease()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	releaseA, err := lease.Acquire(ctx, "BTCUSDT", "A", 10*time.Millisecond)
	assert.NoError(t, err)

	_, err = lease.Acquire(ctx, "BTCUSDT", "B", 10*time.Millisecond)
	assert.Error(t, err)

	assert.NoError(t, releaseA())

	releaseB, err := lease.Acquire(ctx, "BTCUSDT", "B", 10*time.Millisecond)
	assert.NoError(t, err)
	assert.NoError(t, releaseB())
}
