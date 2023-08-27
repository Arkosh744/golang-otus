package memrepo

import (
	"context"
	"testing"
	"time"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/models"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/repo"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	memRepo := NewRepo()
	timeNow := time.Now()

	eventOne := &models.Event{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       "Event title 1",
		StartAt:     timeNow,
		EndAt:       timeNow,
		Description: "Event description",
		UserID:      uuid.Must(uuid.NewV4()),
		NotifyAt:    timeNow,
	}
	eventTwo := &models.Event{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       "Event title 2",
		StartAt:     timeNow.Add(-time.Minute * 15),
		EndAt:       timeNow,
		Description: "Event description",
		UserID:      uuid.Must(uuid.NewV4()),
		NotifyAt:    timeNow,
	}
	eventThree := &models.Event{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       "Event title 3",
		StartAt:     timeNow,
		EndAt:       timeNow.Add(time.Minute * 15),
		Description: "Event description",
		UserID:      uuid.Must(uuid.NewV4()),
		NotifyAt:    timeNow,
	}
	eventFour := &models.Event{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       "Event title 4",
		StartAt:     timeNow,
		EndAt:       timeNow,
		Description: "Event description",
		UserID:      uuid.Must(uuid.NewV4()),
		NotifyAt:    timeNow,
	}

	require.Nil(t, memRepo.CreateEvent(ctx, eventOne))
	require.Nil(t, memRepo.CreateEvent(ctx, eventTwo))
	require.Nil(t, memRepo.CreateEvent(ctx, eventThree))

	resEventOne, err := memRepo.GetEventByID(ctx, eventOne.ID)
	require.Nil(t, err)
	require.Equal(t, eventOne, resEventOne)

	eventsList, err := memRepo.ListEventsByPeriod(ctx, time.Now().Add(-time.Minute*5), time.Now().Add(time.Minute*5), 10)
	require.Nil(t, err)
	require.Equal(t, []*models.Event{eventOne}, eventsList)

	eventsListLimitTwo, err := memRepo.ListEventsByPeriod(ctx, time.Now().Add(-time.Hour*5), time.Now().Add(time.Hour*5), 2)
	require.Nil(t, err)
	require.Equal(t, 2, len(eventsListLimitTwo))

	require.Nil(t, memRepo.DeleteEvent(ctx, eventOne.ID))
	eventOneActual, err := memRepo.GetEventByID(ctx, eventOne.ID)
	require.Nil(t, eventOneActual)
	require.Equal(t, repo.ErrEventNotExist, err)

	eventTwo.Title = "Other title"
	require.Nil(t, memRepo.UpdateEvent(ctx, eventTwo))
	resEventTwoActual, err := memRepo.GetEventByID(ctx, eventTwo.ID)
	require.Nil(t, err)
	require.Equal(t, "Other title", resEventTwoActual.Title)

	require.ErrorIs(t, memRepo.UpdateEvent(ctx, eventFour), repo.ErrEventNotExist)
}
