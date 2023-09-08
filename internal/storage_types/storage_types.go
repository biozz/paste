package storage_types

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Storage interface {
	Init(ctx context.Context) (Closer, error)
	CreatePaste(ctx context.Context, paste Paste) (Paste, error)
	GetPasteBySlug(ctx context.Context, slug string) (Paste, error)
}

type Closer func()

type Paste struct {
	bun.BaseModel `json:"-" bun:"table:pastes,alias:p"`

	ID        int64     `json:"id,omitempty" bun:"id,pk,autoincrement"`
	Slug      string    `json:"slug" bun:"slug,unique"`
	Content   string    `json:"content" bun:"content"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
}
