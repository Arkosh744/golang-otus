package pg

import (
	"context"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/log"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

var _ Client = (*client)(nil)

type Client interface {
	Close() error
	PG() PG
}

type client struct {
	pg PG
}

func NewClient(ctx context.Context, pgCfg *pgxpool.Config) (Client, error) {
	dbc, err := pgxpool.ConnectConfig(ctx, pgCfg)
	if err != nil {
		log.Errorf("failed to connect to postgres", zap.Error(err))

		return nil, err
	}

	log.Info("pg connected successfully")

	return &client{pg: &pg{pgxPool: dbc}}, nil
}

func (c *client) PG() PG {
	return c.pg
}

func (c *client) Close() error {
	if c.pg != nil {
		return c.pg.Close()
	}

	return nil
}
