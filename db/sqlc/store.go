package db

import (
	// "context"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// here the the functions will be written here
// Store defines all functions to execute db queries and transactions

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

var _ Store = &SQLStore{}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)
	q := New(tx)
	qtx := q.WithTx(tx)
	err = fn(qtx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
