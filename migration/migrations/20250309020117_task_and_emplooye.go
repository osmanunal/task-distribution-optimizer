package migrations

import (
	"context"
	"task-distribution-optimizer/migration"
	"task-distribution-optimizer/pkg/model"

	"github.com/uptrace/bun"
)

type Employee struct {
	model.BaseModel

	Name       string `bun:",notnull"`
	Difficulty int    `bun:",notnull"`
	Workload   int    `bun:",notnull"`
}

type Task struct {
	model.BaseModel
	ExternalID int64
	Name       string
	Difficulty int
	Duration   int
	Processed  bool
	EmployeeID int64
	Employee   Employee `bun:"rel:belongs-to,join:employee_id=id"`
}

func init() {
	models := []interface{}{
		&Employee{},
		&Task{},
	}
	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
			if err != nil {
				return err
			}

			for _, m := range models {
				_, err = tx.NewCreateTable().Model(m).IfNotExists().WithForeignKeys().Exec(ctx)
				if err != nil {
					return err
				}
			}

			_, err = tx.NewCreateIndex().Model(&Task{}).Unique().Column("external_id").
				Index("idx_task_external_id").Exec(ctx)
			return nil
		})
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) (err error) {
			for _, m := range models {
				_, err = tx.NewDropTable().Model(m).IfExists().Cascade().Exec(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	migration.Migrations.MustRegister(up, down)
}
