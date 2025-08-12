package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Tx struct {
	db *sql.DB
}

func InitTx(db *sql.DB) *Tx {
	if db == nil {
		panic("db is nil")
	}
	return &Tx{db: db}
}

func (t *Tx) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("tx begin failed: %w", err)
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("tx function failed: %w", err)
	}

	return tx.Commit()
}
