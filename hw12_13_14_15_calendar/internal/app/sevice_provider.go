package app

import (
	"context"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/client/pg"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/closer"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/handlers"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/log"
	memrepo "github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/repo/memory"
	sqlrepo "github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/repo/sql"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type serviceProvider struct {
	calendarService handlers.CalendarService
	pgClient        pg.Client

	repo service.Repository
}

func newServiceProvider(ctx context.Context) *serviceProvider {
	sp := &serviceProvider{}
	sp.GetCalendarService(ctx)

	return sp
}

func (s *serviceProvider) GetPGClient(ctx context.Context) pg.Client {
	if s.pgClient == nil {
		pgCfg, err := pgxpool.ParseConfig(config.AppConfig.GetPostgresDSN())
		if err != nil {
			log.Fatalf("failed to parse pg config", zap.Error(err))
		}

		cl, err := pg.NewClient(ctx, pgCfg)
		if err != nil {
			log.Fatalf("failed to get pg client", zap.Error(err))
		}

		if cl.PG().Ping(ctx) != nil {
			log.Fatalf("failed to ping pg", zap.Error(err))
		}

		closer.Add(cl.Close)

		s.pgClient = cl
	}

	return s.pgClient
}

func (s *serviceProvider) GetCalendarRepo(ctx context.Context) service.Repository {
	if s.repo == nil {
		switch config.AppConfig.GetStorage() {
		case config.StorageMemory:
			s.repo = memrepo.NewRepo()
		case config.StoragePostgres:
			s.repo = sqlrepo.NewRepo(s.GetPGClient(ctx))
		}
	}

	return s.repo
}

func (s *serviceProvider) GetCalendarService(ctx context.Context) handlers.CalendarService {
	if s.calendarService == nil {
		s.calendarService = service.New(s.GetCalendarRepo(ctx))
	}

	return s.calendarService
}
