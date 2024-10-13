package database

import (
	"context"
	"database/sql"
)

type IDatabase interface {
	GetConnection() (*sql.DB, error)
	GetConnectionString() string
	InitDatabase(file string) error
	TryConnection() error
	Migrate() error
}

type IRateLimitDatabase[T any] interface {
	Get(ctx context.Context, key string) (*T, bool)
	Set(ctx context.Context, key string, value *T) error
}
