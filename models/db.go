package models

import (
	"database/sql"
	"fmt"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/sqlite3"
	_ "github.com/mattes/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/user"
	"path"
	"time"
)

func PrepareDB() *sql.DB {
	dir := prepareDBDir()

	dbPath := path.Join(dir, "godestroy.db")
	_, err := os.Stat(dbPath)
	isNew := os.IsNotExist(err)
	db, err := sql.Open("sqlite3", dbPath)
	exitIfError(err)

	prepareMigrations(db, isNew)

	return db
}

func prepareMigrations(db *sql.DB, createTable bool) {
	if createTable {
		err := CreateMigrationTable(db)
		exitIfError(err)
	}

	latestMigration, err := GetLatestMigration(db)
	exitIfError(err)

	if CURRENT_MIGRATION_VERSION == latestMigration.ID {
		return
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	exitIfError(err)

	currentDir, err := os.Getwd()
	exitIfError(err)

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s/migrations", currentDir),
		"sqlite3", driver,
	)
	exitIfError(err)

	step := CURRENT_MIGRATION_VERSION - latestMigration.ID
	// ensure we always migrate up
	if step < 0 {
		step = 1
	}
	err = m.Steps(step)
	exitIfError(err)

	latestMigration.ID = CURRENT_MIGRATION_VERSION
	latestMigration.Time = time.Now()
	err = latestMigration.Write(db)
	exitIfError(err)
}

func prepareDBDir() string {
	currentUser, err := user.Current()
	exitIfError(err)
	dir := path.Join(currentUser.HomeDir, ".godestroy")

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		exitIfError(err)
	}

	return dir
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
