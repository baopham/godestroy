package cli

import (
	"database/sql"
	"github.com/baopham/godestroy/models"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
	"time"
)

func Files(c *cli.Context, db *sql.DB) error {
	var schedules []*models.Schedule

	// TODO: parse --in and --at options
	timeToDestroy := time.Now().AddDate(0, 0, 1)

	for _, path := range c.Args() {
		info, err := os.Stat(path)
		if os.IsNotExist(err) || info.IsDir() {
			color.Yellow("%s is not a valid file", path)
			continue
		}
		schedule := &models.Schedule{
			Path:          path,
			Size:          info.Size(),
			ModTime:       info.ModTime(),
			TimeToDestroy: timeToDestroy,
		}
		schedules = append(schedules, schedule)
	}

	color.Green("Scheduling to destroy %d files...", len(schedules))
	err := models.WriteSchedules(schedules, db)
	color.Green("Done")
	return err
}
