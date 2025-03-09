package database

import (
	"fmt"
	"task-distribution-optimizer/pkg/config"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/schema"
)

func ConnectDB(config config.DBConfig) *bun.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=prefer",
		config.Host, config.Port, config.Name, config.User, config.Password,
	)
	c, err := pgx.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	c.PreferSimpleProtocol = true
	sqlDB := stdlib.OpenDB(*c)
	db := bun.NewDB(sqlDB, pgdialect.New())
	if err = db.Ping(); err != nil {
		panic(err)
	}

	schema.SetTableNameInflector(func(s string) string {
		return s
	})

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return db
}
