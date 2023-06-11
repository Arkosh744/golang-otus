package memrepo

import "sync"

type storage struct {
	mu sync.RWMutex //nolint:unused
}

func NewRepo() *storage {
	return &storage{}
}
