package uow

import (
	"database/sql"
	"fmt"
)

type Transaction[T any] struct {
	db *sql.DB
	fn func() (T, error)
}

func NewTransaction[T any](db *sql.DB, fn func() (T, error)) *Transaction[T] {
	return &Transaction[T]{
		db: db,
		fn: fn,
	}
}
func (t *Transaction[T]) Exec() (T, error) {
	var result T
	tx, err := t.db.Begin()
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err = t.fn()
	if err != nil {
		return result, fmt.Errorf("failed to execute transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, err
}
