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
			Name:   "files",
			Usage:  "godestroy files ~/Desktop/Screen*.png --in 10days",
			Action: action(destroyCli.Files),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "in",
					Usage: "Time to wait to destroy the files (e.g. 2sec, 2min, 2hours, etc.)",
				},
				cli.StringFlag{
					Name:  "at",
					Usage: "Specific time when to destroy the files",
				},
			},
		},
	}

	app.Run(os.Args)
}
