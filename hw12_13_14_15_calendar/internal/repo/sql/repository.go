package sqlrepo

import (
	"context"
	"time"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/client/pg"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

type repository struct {
	client pg.Client
}

func NewRepo(client pg.Client) *repository {
	return &repository{client: client}
}

const tableEvents = "events"

func (r *repository) CreateEvent(ctx context.Context, event *models.Event) error {
	builder := sq.Insert(tableEvents).
		Columns("id", "title", "description", "start", "end", "owner_id").
		Values(event.ID, event.Title, event.Description, event.StartAt, event.EndAt, event.UserID)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "repo.CreateEvent",
		QueryRaw: query,
	}

	if _, err = r.client.PG().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateEvent(ctx context.Context, event *models.Event) error {
	builder := sq.Update(tableEvents).
		Set("title", event.Title).
		Set("description", event.Description).
		Set("start", event.StartAt).
		Set("end", event.EndAt).
		Where(sq.Eq{"id": event.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "repo.UpdateEvent",
		QueryRaw: query,
	}

	if _, err = r.client.PG().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	builder := sq.Delete(tableEvents).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "repo.DeleteEvent",
		QueryRaw: query,
	}

	if _, err = r.client.PG().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	builder := sq.Select("id", "title", "description", "start_at", "end_at").
		From(tableEvents).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "repo.GetEventByID",
		QueryRaw: query,
	}

	var event *models.Event
	if err := r.client.PG().ScanOneContext(ctx, &event, q, args...); err != nil {
		return nil, err
	}

	return event, nil
}

func (r *repository) ListEventsByPeriod(ctx context.Context, start, end time.Time, limit int) ([]*models.Event, error) {
	builder := sq.Select("id", "title", "description", "start_at", "end_at").
		From(tableEvents).
		Where(sq.And{
			sq.GtOrEq{"start_at": start},
			sq.LtOrEq{"end_at": end},
		}).
		Limit(uint64(limit))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "repo.ListEventsByPeriod",
		QueryRaw: query,
	}

	var res []*models.Event
	if err := r.client.PG().ScanAllContext(ctx, &res, q, args...); err != nil {
		return nil, err
	}

	return res, nil
}
