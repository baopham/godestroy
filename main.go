package main

import (
	"database/sql"
	destroyCli "github.com/baopham/godestroy/cli"
	"github.com/baopham/godestroy/models"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
)

func main() {
	db := models.PrepareDB()
	defer db.Close()

	action := func(fn func(c *cli.Context, db *sql.DB) error) func(c *cli.Context) error {
		actor := func(c *cli.Context) error {
			err := fn(c, db)
			if err != nil {
				color.Red(err.Error())
			}
			return err
		}
		return actor
	}

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "godestroy"
	app.Usage = "Schedule to destroy file(s)"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:        "files",
			Usage:       "godestroy files ~/Desktop/Screen*.png --in 10days",
			Description: "Schedule to destroy the provided files",
			Action:      action(destroyCli.Schedule),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "in",
					Usage: "Time to wait to destroy the files (e.g. --in 2seconds, 2mins, 2hours, etc.)",
				},
				cli.StringFlag{
					Name:  "at",
					Usage: `Specific time when to destroy the files (e.g. --at "November 11, 2020")`,
				},
			},
		},
		{
			Name:        "what?",
			Usage:       "godestroy what?",
			Description: "List all the scheduled files",
			Aliases:     []string{"list"},
			Action:      action(destroyCli.List),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "now",
					Usage: "List files that should be destroyed now",
				},
			},
		},
		{
			Name:        "not",
			Usage:       "godestroy not files ~/Desktop/Screen*.png",
			Description: "Don't destroy the provided files",
			Aliases:     []string{"remove"},
			Action:      action(destroyCli.Remove),
		},
		{
			Name:        "now!",
			Usage:       "godestroy now!",
			Description: "Destroy the files that are scheduled to be deleted now",
			Aliases:     []string{"remove"},
			Action:      action(destroyCli.Destroy),
		},
	}

	app.Run(os.Args)
}
