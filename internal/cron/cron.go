package scheduler

import (
	"context"
	"log/slog"

	goCron "github.com/robfig/cron/v3"
)

type scheduler struct {
	logger *slog.Logger
	cron   *goCron.Cron
}

func New(logger *slog.Logger, cron *goCron.Cron) *scheduler {
	return &scheduler{
		logger: logger,
		cron:   cron,
	}
}

func (s *scheduler) AddFunc(ctx context.Context, spec string, cmd func()) error {
	s.logger.DebugContext(ctx, "Adding cron job", "spec", spec)
	entry, err := s.cron.AddFunc(spec, cmd)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to add cron job", "spec", spec, "error", err)
		return err
	}
	s.logger.DebugContext(ctx, "Cron job added", "spec", spec, "entryID", entry)
	return nil
}

func (s *scheduler) Start(ctx context.Context) error {
	s.logger.DebugContext(ctx, "Starting cron scheduler")
	s.cron.Start()
	return nil
}

func (s *scheduler) Stop(ctx context.Context) error {
	s.logger.DebugContext(ctx, "Stopping cron scheduler")
	s.cron.Stop()
	return nil
}
