package sqlrepo

import (
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/client/pg"
)

type storage struct {
	client pg.Client
}

func NewRepo(client pg.Client) *storage {
	return &storage{client: client}
}
