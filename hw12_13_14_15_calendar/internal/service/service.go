package service

import (
	"context"
	"time"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/models"
	"github.com/gofrs/uuid"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	ListEventsByPeriod(ctx context.Context, start, end time.Time, limit int) ([]*models.Event, error)
}

type calendarService struct {
	repo Repository
}

func New(repo Repository) *calendarService {
	return &calendarService{
		repo: repo,
	}
}
