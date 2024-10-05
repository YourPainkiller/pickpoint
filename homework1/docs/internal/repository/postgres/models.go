package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryEngine interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)

	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)

	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
}

type TransactionManager interface {
	GetQueryEngine(ctx context.Context) QueryEngine
	RunSerialazible(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunReadUncommited(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunReadWriteCommited(ctx context.Context, fn func(ctxTx context.Context) error) error
}
