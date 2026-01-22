package matching

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type JetStreamWAL struct {
	js nats.JetStreamContext
}

func NewJetStreamWAL(js nats.JetStreamContext) (*JetStreamWAL, error) {
	if err := ensureMatchingEventStream(js); err != nil {
		return nil, err
	}
	return &JetStreamWAL{js: js}, nil
}

func ensureMatchingEventStream(js nats.JetStreamContext) error {
	_, err := js.StreamInfo("MATCHING")
	if err == nil {
		return nil
	}
	if err != nats.ErrStreamNotFound {
		return err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "MATCHING",
		Subjects:  []string{"matching.event.*"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxAge:    7 * 24 * time.Hour,
	})
	return err
}

func (w *JetStreamWAL) Append(ev Event) error {
	subject := fmt.Sprintf("matching.event.%s", ev.Symbol)
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	_, err = w.js.Publish(subject, data)
	return err
}

func (w *JetStreamWAL) Load(symbol string, afterSeq uint64) ([]Event, error) {
	subject := fmt.Sprintf("matching.event.%s", symbol)
	sub, err := w.js.PullSubscribe(subject, "", nats.BindStream("MATCHING"))
	if err != nil {
		return nil, err
	}
	defer func() { _ = sub.Unsubscribe() }()

	var out []Event
	for {
		msgs, err := sub.Fetch(256, nats.MaxWait(100*time.Millisecond))
		if err != nil {
			if err == nats.ErrTimeout {
				break
			}
			return nil, err
		}
		for _, msg := range msgs {
			var ev Event
			if err := json.Unmarshal(msg.Data, &ev); err == nil {
				if ev.Seq > afterSeq {
					out = append(out, ev)
				}
			}
			_ = msg.Ack()
		}
	}
	return out, nil
}

type JetStreamSnapshotStore struct {
	kv nats.KeyValue
}

func NewJetStreamSnapshotStore(js nats.JetStreamContext) (*JetStreamSnapshotStore, error) {
	kv, err := js.KeyValue("MATCHING_SNAPSHOT")
	if err == nats.ErrBucketNotFound {
		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:  "MATCHING_SNAPSHOT",
			History: 1,
			Storage: nats.FileStorage,
		})
	}
	if err != nil {
		return nil, err
	}
	return &JetStreamSnapshotStore{kv: kv}, nil
}

func (s *JetStreamSnapshotStore) Put(symbol string, snap EngineSnapshot) error {
	data, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	_, err = s.kv.Put(symbol, data)
	return err
}

func (s *JetStreamSnapshotStore) Get(symbol string) (*EngineSnapshot, error) {
	entry, err := s.kv.Get(symbol)
	if err != nil {
		if err == nats.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	var snap EngineSnapshot
	if err := json.Unmarshal(entry.Value(), &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

type JetStreamLease struct {
	kv nats.KeyValue
}

func NewJetStreamLease(js nats.JetStreamContext, ttl time.Duration) (*JetStreamLease, error) {
	kv, err := js.KeyValue("MATCHING_LEASE")
	if err == nats.ErrBucketNotFound {
		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:  "MATCHING_LEASE",
			History: 1,
			Storage: nats.FileStorage,
			TTL:     ttl,
		})
	}
	if err != nil {
		return nil, err
	}
	return &JetStreamLease{kv: kv}, nil
}

func (l *JetStreamLease) Acquire(ctx context.Context, symbol string, owner string, refresh time.Duration) (func() error, error) {
	_, err := l.kv.Create(symbol, []byte(owner))
	if err != nil {
		return nil, err
	}
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(refresh)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_, _ = l.kv.Put(symbol, []byte(owner))
			case <-stop:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
	return func() error {
		close(stop)
		return l.kv.Delete(symbol)
	}, nil
}
