package models

import (
	"time"
)

type Notification struct {
	ID       string
	Title    string
	Datetime time.Time
	UserTo   UserID
}
