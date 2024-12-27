package service

import (
	"context"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
	"log/slog"
	"time"
)

type Statistics interface {
	GetTotalUsers() (int, error)
	GetTotalConnections() (int, error)

	GetUserConnections(ctx context.Context, handle entity.UserHandle) (int, error)

	do.Shutdownable
}

type statisticsimpl struct {
	logger                *slog.Logger
	userRepo              repository.User
	totalUsersCache       int
	totalConnectionsCache int
	updateErr             error
	ticker                *time.Ticker
	updateTimeout         time.Duration
}

func NewStatisticsService(injector *do.Injector) (Statistics, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}

	userRepo, err := do.Invoke[repository.User](injector)
	if err != nil {
		return nil, fmt.Errorf("getting user repository: %w", err)
	}

	tickerUpdate, err := time.ParseDuration(config.Statistics.UpdateInterval)
	if err != nil {
		return nil, fmt.Errorf("parsing update interval: %w", err)
	}

	tickerUpdateTimeout, err := time.ParseDuration(config.Statistics.UpdateTimeout)
	if err != nil {
		return nil, fmt.Errorf("parsing update interval: %w", err)
	}

	stats := &statisticsimpl{
		logger:                logger.With("service", "statistics"),
		userRepo:              userRepo,
		totalUsersCache:       0,
		totalConnectionsCache: 0,
		ticker:                time.NewTicker(tickerUpdate),
		updateTimeout:         tickerUpdateTimeout,
	}
	stats.startStatsCacheUpdate()

	return stats, nil
}

func (s *statisticsimpl) startStatsCacheUpdate() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
		defer cancel()

		s.updateErr = s.updateStatsCache(ctx)
		if s.updateErr != nil {
			s.logger.Error("failed to update statistics", "error", s.updateErr.Error())
		}

		for {
			_, ok := <-s.ticker.C
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
			defer cancel()

			s.updateErr = s.updateStatsCache(ctx)
			if s.updateErr != nil {
				s.logger.Error("failed to update statistics", "error", s.updateErr.Error())
			}
		}
	}()
}

func (s *statisticsimpl) updateStatsCache(ctx context.Context) error {
	var err error
	defer s.logger.Debug("updated statistics cache", "totalUsers", s.totalUsersCache, "totalConnections", s.totalConnectionsCache, "updateErr", err)

	s.totalConnectionsCache, err = s.userRepo.GetOverallConnectionsCount(ctx)
	if err != nil {
		return fmt.Errorf("getting total connections count: %w", err)
	}

	s.totalUsersCache, err = s.userRepo.GetOverallUserCount(ctx)
	if err != nil {
		return fmt.Errorf("getting total users count: %w", err)
	}

	return nil
}

func (s *statisticsimpl) GetTotalUsers() (int, error) {
	if s.updateErr != nil {
		return 0, s.updateErr
	}
	return s.totalUsersCache, nil
}

func (s *statisticsimpl) GetTotalConnections() (int, error) {
	if s.updateErr != nil {
		return 0, s.updateErr
	}
	return s.totalConnectionsCache, nil
}

func (s *statisticsimpl) GetUserConnections(ctx context.Context, handle entity.UserHandle) (int, error) {
	count, err := s.userRepo.GetUserConnectionsCount(ctx, handle)
	if err != nil {
		return 0, fmt.Errorf("getting user connections count: %w", err)
	}
	return count, nil
}

func (s *statisticsimpl) Shutdown() error {
	s.ticker.Stop()
	return nil
}
