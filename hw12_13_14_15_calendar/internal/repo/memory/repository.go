package memrepo

import (
	"context"
	"sync"
	"time"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/models"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/repo"
	"github.com/gofrs/uuid"
)

type storage struct {
	events map[uuid.UUID]*models.Event
	mu     sync.RWMutex
}

func NewRepo() *storage {
	return &storage{
		events: make(map[uuid.UUID]*models.Event),
	}
}

func (m *storage) CreateEvent(_ context.Context, e *models.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events[e.ID] = e

	return nil
}

func (m *storage) UpdateEvent(_ context.Context, e *models.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events[e.ID] = e

	return nil
}

func (m *storage) DeleteEvent(_ context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.events, id)

	return nil
}

func (m *storage) GetEventByID(_ context.Context, id uuid.UUID) (*models.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	event, ok := m.events[id]
	if !ok {
		return nil, repo.ErrEventNotExist
	}

	return event, nil
}

func (m *storage) ListEventsByPeriod(_ context.Context, start, end time.Time, limit int) ([]*models.Event, error) {
	var res []*models.Event
	for _, e := range m.events {
		if len(res) >= limit {
			break
		}

		if e.StartAt.After(start) && e.EndAt.Before(end) {
			res = append(res, e)
		}
	}
	return res, nil
}
