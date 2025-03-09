#!/bin/bash

MIGRATION_DIR="migration/migrations"
mkdir -p $MIGRATION_DIR

TIMESTAMP=$(date +%Y%m%d%H%M%S)
DESCRIPTION=$(echo $1 | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
FILENAME="${MIGRATION_DIR}/${TIMESTAMP}_${DESCRIPTION}.go"

# Create migration file
cat > $FILENAME << EOF
package migrations

import (
	"context"
	"taskmanager/migration"

	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) (err error) {
			return err
		})
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) (err error) {
			return err
		})
	}
	migration.Migrations.MustRegister(up, down)
}
EOF

echo "Yeni migration dosyası oluşturuldu: $FILENAME" 