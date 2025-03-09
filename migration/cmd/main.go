package main

import (
	"context"
	"os"
	"taskmanager/migration"
	_ "taskmanager/migration/migrations"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "bun",
		Commands: commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

var commands = []*cli.Command{
	{
		Name:  "migrate",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "up",
				Usage: "run all available migrations",
				Action: func(c *cli.Context) error {
					return migration.Up(context.Background())
				},
			},
			{
				Name:  "down",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					return migration.Down(context.Background())
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					return migration.Status(context.Background())
				},
			},
		},
	},
}
