package cli

import (
	"database/sql"
	"github.com/baopham/godestroy/models"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
)

func Remove(c *cli.Context, db *sql.DB) error {
	for _, path := range c.Args()[1:] {
		info, err := os.Stat(path)
		if os.IsNotExist(err) || info.IsDir() {
			color.Yellow("%s is not a valid file", path)
			continue
		}
		schedule, err := models.FindSchedule(path, db)
		if err != nil {
			if err == sql.ErrNoRows {
				// be forgiving...
				color.Yellow("%s has not been scheduled", path)
				continue
			}
			return err
		}
		err = schedule.Remove(db)
		if err != nil {
			return err
		}
		color.Green("%s has been descheduled", path)
	}
	return nil
}
