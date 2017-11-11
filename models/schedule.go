package models

import (
	"database/sql"
	"time"
)

type Schedule struct {
	ID            int
	Path          string
	Size          int64
	ModTime       time.Time
	TimeToDestroy time.Time
}

func WriteSchedules(schedules []*Schedule, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, schedule := range schedules {
		_, err = tx.Exec(`
			insert into schedules(path, size, mod_time, time_to_destroy)
			values(?, ?, ?, ?)
		`, schedule.Path, schedule.Size, schedule.ModTime, schedule.TimeToDestroy)
	}

	return tx.Commit()
}

func GetAllSchedules(db *sql.DB) ([]*Schedule, error) {
	var schedules []*Schedule

	rows, err := db.Query("select id, path, size, mod_time, time_to_destroy from schedules")
	if err != nil {
		return schedules, err
	}

	for rows.Next() {
		var id int
		var path string
		var size int64
		var modTime time.Time
		var timeToDestroy time.Time
		err := rows.Scan(&id, &path, &size, &modTime, &timeToDestroy)
		if err != nil {
			return schedules, err
		}
		schedule := &Schedule{
			ID:            id,
			Path:          path,
			Size:          size,
			ModTime:       modTime,
			TimeToDestroy: timeToDestroy,
		}
		schedules = append(schedules, schedule)
	}

	return schedules, rows.Close()
}
