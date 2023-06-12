package memrepo

import (
	"context"
	"sync"
	"time"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/models"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/repo"
	"github.com/gofrs/uuid"
)

type repository struct {
	events map[uuid.UUID]*models.Event
	mu     sync.Mutex
}

func NewRepo() *repository {
	return &repository{
		events: make(map[uuid.UUID]*models.Event),
	}
}

func (r *repository) CreateEvent(_ context.Context, event *models.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[event.ID] = event

	return nil
}

func (r *repository) UpdateEvent(_ context.Context, event *models.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.events[event.ID]
	if !ok {
		return repo.ErrEventNotExist
	}

	r.events[event.ID] = event

	return nil
}

func (r *repository) DeleteEvent(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.events, id)

	return nil
}

func (r *repository) GetEventByID(_ context.Context, id uuid.UUID) (*models.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, ok := r.events[id]
	if !ok {
		return nil, repo.ErrEventNotExist
	}

	return event, nil
}

func (r *repository) ListEventsByPeriod(_ context.Context, start, end time.Time, limit int) ([]*models.Event, error) {
	var res []*models.Event

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, e := range r.events {
		if len(res) >= limit {
			break
		}

		if e.StartAt.After(start) && e.EndAt.Before(end) {
			res = append(res, e)
		}
	}

	return res, nil
}
