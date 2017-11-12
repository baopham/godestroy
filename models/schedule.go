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

func (s *Schedule) Remove(db *sql.DB) error {
	_, err := db.Exec("delete from schedules where id = ?", s.ID)
	return err
}

func FindSchedule(path string, db *sql.DB) (*Schedule, error) {
	row := db.QueryRow("select id, path, size, mod_time, time_to_destroy from schedules where path = ?", path)
	schedule := &Schedule{}
	err := row.Scan(&schedule.ID, &schedule.Path, &schedule.Size, &schedule.ModTime, &schedule.TimeToDestroy)
	if err != nil {
		return nil, err
	}
	return schedule, nil
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

func GetDueSchedules(db *sql.DB) ([]*Schedule, error) {
	rows, err := db.Query(
		"select id, path, size, mod_time, time_to_destroy from schedules where time_to_destroy <= ?",
		time.Now(),
	)
	if err != nil {
		return []*Schedule{}, err
	}

	return rowsToSchedules(rows)
}

func GetAllSchedules(db *sql.DB) ([]*Schedule, error) {
	rows, err := db.Query("select id, path, size, mod_time, time_to_destroy from schedules")
	if err != nil {
		return []*Schedule{}, err
	}

	return rowsToSchedules(rows)
}

func rowsToSchedules(rows *sql.Rows) ([]*Schedule, error) {
	var schedules []*Schedule

	for rows.Next() {
		schedule := &Schedule{}
		err := rows.Scan(&schedule.ID, &schedule.Path, &schedule.Size, &schedule.ModTime, &schedule.TimeToDestroy)
		if err != nil {
			return schedules, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, rows.Close()
}
