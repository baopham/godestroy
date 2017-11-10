package cli

import (
	"database/sql"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
)

func Files(c *cli.Context, db *sql.DB) error {
	for _, path := range c.Args() {
		info, err := os.Stat(path)
		if os.IsNotExist(err) || info.IsDir() {
			color.Yellow("%s is not a valid file", path)
			continue
		}
	}

	return nil
}
