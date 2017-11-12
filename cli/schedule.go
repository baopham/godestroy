package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/baopham/godestroy/models"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"time"
)

func Schedule(c *cli.Context, db *sql.DB) error {
	timeToDestroy, err := parseTimeOption(c)
	if err != nil {
		return err
	}

	var schedules []*models.Schedule
	for _, path := range c.Args() {
		valid, info, err := isValidFile(path)
		if !valid {
			color.Yellow("%s is not a valid file", path)
			continue
		}
		if err != nil {
			return err
		}
		if schedule, _ := models.FindSchedule(path, db); schedule != nil {
			color.Yellow("%s is already scheduled to be destroyed", path)
			continue
		}
		schedule := &models.Schedule{
			Path:          path,
			Size:          info.Size(),
			ModTime:       info.ModTime(),
			TimeToDestroy: *timeToDestroy,
		}
		schedules = append(schedules, schedule)
	}

	color.Green("Scheduling to destroy %d files...", len(schedules))
	err = models.WriteSchedules(schedules, db)
	color.Green("Done")
	return err
}

func parseTimeOption(c *cli.Context) (*time.Time, error) {
	tParser := getTimeParser()

	if in := c.String("in"); in != "" {
		timeToDestroy, err := tParser.Parse("in "+in, time.Now())
		if timeToDestroy == nil {
			return nil, errors.New(fmt.Sprintf("Could not parse %s", in))
		}
		return &timeToDestroy.Time, err
	}

	if at := c.String("at"); at != "" {
		timeToDestroy, err := tParser.Parse("at "+at, time.Now())
		if timeToDestroy == nil {
			return nil, errors.New(fmt.Sprintf("Could not parse %s", at))
		}
		return &timeToDestroy.Time, err
	}

	return nil, errors.New("No time provided")
}
