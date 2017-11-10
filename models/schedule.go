package models

import (
	"time"
)

const SCHEDULE_TABLE = "schedules"

type Schedule struct {
	ID            int
	Path          string
	Size          int64
	ModTime       time.Time
	TimeToDestroy time.Time
}
