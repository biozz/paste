package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/biozz/paste/internal/storage/migrations"
	"github.com/biozz/paste/internal/storage_types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	_ "modernc.org/sqlite"
)

type SQLStorage struct {
	// inputs
	uri string
	dev bool

	// internal
	db *bun.DB
}

func NewSQL(uri string) *SQLStorage {
	return &SQLStorage{
		uri: uri,
		dev: true,
	}
}

func (s *SQLStorage) Init(ctx context.Context) (storage_types.Closer, error) {
	if strings.HasPrefix(s.uri, "file:") {
		filename := strings.TrimLeft(s.uri, "file:")
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			_, err = os.Create(filename)
			if err != nil {
				return nil, err
			}
		}
	}
	if strings.HasPrefix(s.uri, "postgres://") {
		sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(s.uri)))
		s.db = bun.NewDB(sqlDB, pgdialect.New())
	} else {
		sqlDB, err := sql.Open("sqlite", s.uri)
		if err != nil {
			return nil, err
		}
		s.db = bun.NewDB(sqlDB, sqlitedialect.New())
	}
	if s.dev {
		s.db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	err := s.migrate(ctx)
	if err != nil {
		return nil, err
	}
	return func() {}, nil
}

func (s *SQLStorage) migrate(ctx context.Context) error {
	migrator := migrate.NewMigrator(s.db, migrations.Migrations)
	err := migrator.Init(ctx)
	if err != nil {
		return err
	}
	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx) //nolint:errcheck
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		fmt.Printf("there are no new migrations to run (database is up to date)\n")
		return nil
	}
	fmt.Printf("migrated to %s\n", group)
	return nil
}

func (s *SQLStorage) CreatePaste(ctx context.Context, paste storage_types.Paste) (storage_types.Paste, error) {
	_, err := s.db.NewInsert().Model(&paste).Exec(ctx)
	if err != nil {
		return paste, err
	}
	return paste, err

}

func (s *SQLStorage) GetPasteBySlug(ctx context.Context, slug string) (storage_types.Paste, error) {
	paste := new(storage_types.Paste)
	if err := s.db.NewSelect().Model(paste).Where("slug = ?", slug).Scan(ctx); err != nil {
		return storage_types.Paste{}, err
	}
	return *paste, nil

}
