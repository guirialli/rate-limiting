package database

import "database/sql"

type IDatabase interface {
	GetConnection() (*sql.DB, error)
	GetConnectionString() string
	InitDatabase(file string) error
	TryConnection() error
	Migrate() error
}
