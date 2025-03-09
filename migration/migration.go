package migration

import (
	"context"
	"fmt"
	"taskmanager/pkg/config"
	"taskmanager/pkg/database"

	"github.com/uptrace/bun/migrate"
)

// Migrations holds all migration operations
var Migrations = migrate.NewMigrations()

func Migrator(ctx context.Context) (*migrate.Migrator, error) {
	cfg := config.Read()

	db := database.ConnectDB(cfg.DBConfig)
	return migrate.NewMigrator(db, Migrations), nil
}

// Up runs all available migrations
func Up(ctx context.Context) error {
	migrator, err := Migrator(ctx)
	if err != nil {
		return err
	}

	if err = migrator.Init(ctx); err != nil {
		return err
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.IsZero() {
		fmt.Printf("there are no new migrations to run\n")
		return nil
	}

	fmt.Printf("migrated to %s\n", group)
	return nil
}

// Down rolls back the last migration group
func Down(ctx context.Context) error {
	migrator, err := Migrator(ctx)
	if err != nil {
		return err
	}

	if err = migrator.Init(ctx); err != nil {
		return err
	}

	group, err := migrator.Rollback(ctx)
	if err != nil {
		return err
	}

	if group.IsZero() {
		fmt.Printf("there are no groups to roll back\n")
		return nil
	}

	fmt.Printf("rolled back %s\n", group)
	return nil
}

// Status prints the status of migrations
func Status(ctx context.Context) error {
	migrator, err := Migrator(ctx)
	if err != nil {
		return err
	}

	if err = migrator.Init(ctx); err != nil {
		return err
	}

	ms, err := migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("migrations: %s\n", ms)
	fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
	fmt.Printf("last migration group: %s\n", ms.LastGroup())

	return nil
}
