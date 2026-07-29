package scheduler

import (
	"errors"
	"sync"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
)

type Scheduler struct {
	interval time.Duration
	log      *logger.Logger

	latest  metrics.Info
	lastErr error

	mu     sync.RWMutex
	ticker *time.Ticker
	done   chan struct{}
	wg     sync.WaitGroup
	once   sync.Once
}

func New(interval time.Duration, log *logger.Logger) *Scheduler {
	return &Scheduler{
		interval: interval,
		log:      log,
	}
}

func (s *Scheduler) Start() error {
	snapshot, err := metrics.GetInfo()
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.latest = snapshot
	s.lastErr = nil
	s.mu.Unlock()

	s.ticker = time.NewTicker(s.interval)
	s.done = make(chan struct{})
	s.wg.Add(1)

	go s.run()

	s.log.Info("metrics scheduler started, interval: %s", s.interval)

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
	snapshot, err := metrics.GetInfo()
	if err != nil {
		s.log.Error("failed to collect metrics: %v", err)

		s.mu.Lock()
		s.lastErr = err
		s.mu.Unlock()

		return
	}

	s.mu.Lock()
	s.latest = snapshot
	s.lastErr = nil
	s.mu.Unlock()
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
	})
}
