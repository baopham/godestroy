package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	CURRENT_MIGRATION_VERSION = 1
)

type Migration struct {
	ID   int
	Time time.Time
}

func GetLatestMigration(db *sql.DB) (*Migration, error) {
	migration := &Migration{}
	stm := fmt.Sprintf(`select id, time from migrations order by time desc`)
	err := db.QueryRow(stm).Scan(&migration.ID, &migration.Time)

	if err == sql.ErrNoRows {
		return migration, nil
	}

	if err != nil {
		return migration, err
	}

	return migration, nil
}

func CreateMigrationTable(db *sql.DB) error {
	stm := fmt.Sprintf(`create table migrations (id integer not null primary key, time datetime)`)
	_, err := db.Exec(stm)
	return err
}

func (m *Migration) Write(db *sql.DB) error {
	_, err := db.Exec(
		`insert into migrations(id, time) values(:id, :time)`,
		sql.Named("id", m.ID),
		sql.Named("time", m.Time),
	)
	return err
}
