package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	NotifyAt    time.Time `db:"notify_at"`

	CreatedAy time.Time `db:"created_at"`
}
