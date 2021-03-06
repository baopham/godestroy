package cli

import (
	"database/sql"
	"github.com/baopham/godestroy/models"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"os"
)

func List(c *cli.Context, db *sql.DB) error {
	fn := models.GetAllSchedules

	if c.Bool("now") {
		fn = models.GetDueSchedules
	}

	schedules, err := fn(db)

	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path", "Time To Destroy"})
	for _, schedule := range schedules {
		table.Append([]string{schedule.Path, schedule.TimeToDestroy.String()})
	}
	table.Render()
	return nil
}
