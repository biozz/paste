package storage

import (
	"errors"
	"strings"

	"github.com/biozz/paste/internal/storage_types"
)

var (
	errUnsupportedURI = errors.New("unsupported DB_DSN structure, file:..., :memory: or postgres://...")
)

func New(dsn string) (storage_types.Storage, error) {
	var s storage_types.Storage
	if strings.HasPrefix(dsn, "file:") || dsn == ":memory:" || strings.HasPrefix(dsn, "postgres://") {
		s = NewSQL(dsn)
	} else {
		return nil, errUnsupportedURI
	}
	return s, nil
}
