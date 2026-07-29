package scheduler

import (
	"errors"
	"sync"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/stream"
)

type Scheduler struct {
	interval time.Duration
	log      *logger.Logger
	store    *sqlite.Store
	hub      *stream.Hub

	latest  metrics.Info
	lastErr error

	stats Stats

	mu     sync.RWMutex
	ticker *time.Ticker
	done   chan struct{}
	wg     sync.WaitGroup
	once   sync.Once
}

type Stats struct {
	StartedAt              time.Time
	LastCollectionAt       time.Time
	LastCollectionDuration time.Duration

	SuccessfulCollections uint64
	FailedCollections     uint64

	LastError string
}

func (s *Scheduler) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stats
}

func New(interval time.Duration, log *logger.Logger, store *sqlite.Store, hub *stream.Hub) *Scheduler {
	return &Scheduler{
		interval: interval,
		log:      log,
		store:    store,
		hub:      hub,
	}
}

func (s *Scheduler) Start() error {
	start := time.Now()
	snapshot, err := metrics.GetInfo()
	duration := time.Since(start)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.latest = snapshot
	s.lastErr = nil
	s.stats.StartedAt = time.Now()
	s.stats.LastCollectionAt = s.stats.StartedAt
	s.stats.LastCollectionDuration = duration
	s.stats.SuccessfulCollections++
	s.stats.LastError = ""
	s.mu.Unlock()

	if s.store != nil {
		if err := s.store.Save(snapshot); err != nil {
			s.log.Error("failed to save initial metrics snapshot: %v", err)
		}
	}

	if s.hub != nil {
		s.hub.Broadcast(snapshot)
	}

	s.ticker = time.NewTicker(s.interval)
	s.done = make(chan struct{})
	s.wg.Add(1)

	go s.run()

	s.log.Info("background scheduler started, interval: %s", s.interval)

	return nil
}

func (s *Scheduler) run() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ticker.C:
			s.collect()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) collect() {
	start := time.Now()
	snapshot, err := metrics.GetInfo()
	duration := time.Since(start)

	s.mu.Lock()

	if err != nil {
		s.log.Error("failed to collect metrics: %v", err)

		s.stats.FailedCollections++
		s.stats.LastError = err.Error()
		s.lastErr = err

		s.mu.Unlock()

		return
	}

	s.latest = snapshot
	s.lastErr = nil

	s.stats.SuccessfulCollections++
	s.stats.LastCollectionAt = time.Now()
	s.stats.LastCollectionDuration = duration
	s.stats.LastError = ""

	s.mu.Unlock()

	if s.store != nil {
		if err := s.store.Save(snapshot); err != nil {
			s.log.Error("failed to save metrics snapshot: %v", err)
		}
	}

	if s.hub != nil {
		s.hub.Broadcast(snapshot)
	}
}

func (s *Scheduler) Snapshot() (metrics.Info, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.latest.Metadata.Timestamp.IsZero() {
		return metrics.Info{}, errors.New("no snapshot available")
	}

	return s.latest, nil
}

func (s *Scheduler) Stop() {
	s.once.Do(func() {
		if s.done != nil {
			close(s.done)
		}
		if s.ticker != nil {
			s.ticker.Stop()
		}
		s.wg.Wait()

		s.log.Info("background scheduler stopped")
	})
}
