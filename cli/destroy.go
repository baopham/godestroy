package cli

import (
	"database/sql"
	"github.com/baopham/godestroy/models"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
)

func Destroy(c *cli.Context, db *sql.DB) error {
	schedules, err := models.GetDueSchedules(db)
	if err != nil {
		return err
	}
	for _, schedule := range schedules {
		err := moveToTrash(schedule)
		if err != nil {
			color.Yellow("Failed to delete %s", schedule.Path)
		}
		color.Green("%s is deleted", schedule.Path)
	}
	return nil
}

func moveToTrash(schedule *models.Schedule) error {
	path := schedule.Path
	valid, info, err := isValidFile(path)
	if !valid {
		color.Yellow("%s is not a valid file anymore. Cannot delete it", path)
		return nil
	}
	if err != nil {
		return err
	}
	if info.Size() != schedule.Size {
		color.Yellow("Size of %s has changed. Cannot delete it", path)
		return nil
	}
	if !info.ModTime().Equal(schedule.ModTime) {
		color.Yellow("%s has been modified since it was last scheduled for deletion. Cannot delete it", path)
		return nil
	}

	// TODO: an option to do a soft delete (e.g. move to Trash folder)
	// TODO: build a report of files being deleted
	return os.Remove(path)
}
