package scheduler

import (
	"context"
	"log/slog"
	"time"
)

type scheduler struct {
	logger  *slog.Logger
	tickers []*time.Ticker
	stop    chan struct{}
}

func New(logger *slog.Logger) *scheduler {
	return &scheduler{
		logger: logger,
		stop:   make(chan struct{}),
	}
}

func (s *scheduler) AddFunc(interval time.Duration, cmd func()) {
	ticker := time.NewTicker(interval)
	s.tickers = append(s.tickers, ticker)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.logger.Debug("Running scheduled task", "interval", interval)
				cmd()
			case <-s.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *scheduler) Start(_ context.Context) {
	s.logger.Debug("Scheduler started")
}

func (s *scheduler) Stop() {
	s.logger.Debug("Stopping scheduler")
	close(s.stop)
}
