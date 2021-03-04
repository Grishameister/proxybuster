package database

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type IDbConn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}
