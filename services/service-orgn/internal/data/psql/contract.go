//go:generate mockgen -source=contract.go -destination=psqlmock/contract.go -package=psqlmock
package psql

import (
	"context"
	"database/sql"
)

type Conns interface {
	Conn(ctx context.Context) (Conn, error)
}

type Conn interface {
	Exec(ctx context.Context, query string, args ...any) (Result, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
}

type Result sql.Result

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type Row interface {
	Scan(dest ...any) error
}
