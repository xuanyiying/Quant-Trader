package matching

import (
	"sync"
)

type WAL interface {
	Append(ev Event) error
	Load(symbol string, afterSeq uint64) ([]Event, error)
}

type MemoryWAL struct {
	mu     sync.RWMutex
	events map[string][]Event
}

func NewMemoryWAL() *MemoryWAL {
	return &MemoryWAL{events: make(map[string][]Event)}
}

func (w *MemoryWAL) Append(ev Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.events[ev.Symbol] = append(w.events[ev.Symbol], ev)
	return nil
}

func (w *MemoryWAL) Load(symbol string, afterSeq uint64) ([]Event, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	src := w.events[symbol]
	if len(src) == 0 {
		return nil, nil
	}
	out := make([]Event, 0, len(src))
	for i := range src {
		if src[i].Seq > afterSeq {
			out = append(out, src[i])
		}
	}
	return out, nil
}

type SnapshotStore interface {
	Put(symbol string, snap EngineSnapshot) error
	Get(symbol string) (*EngineSnapshot, error)
}

type MemorySnapshotStore struct {
	mu        sync.RWMutex
	snapshots map[string]EngineSnapshot
}

func NewMemorySnapshotStore() *MemorySnapshotStore {
	return &MemorySnapshotStore{snapshots: make(map[string]EngineSnapshot)}
}

func (s *MemorySnapshotStore) Put(symbol string, snap EngineSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots[symbol] = snap
	return nil
}

func (s *MemorySnapshotStore) Get(symbol string) (*EngineSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.snapshots[symbol]
	if !ok {
		return nil, nil
	}
	return &snap, nil
}
