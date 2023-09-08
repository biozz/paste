package migrations

import (
	"context"

	"github.com/biozz/paste/internal/storage_types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model((*storage_types.Paste)(nil)).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().Model((*storage_types.Paste)(nil)).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
