package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
)

type trxKey struct{}

type conns struct {
	db     *sql.DB
	logger log.Logger
}

func Connect(ctx context.Context, url string, logger log.Logger) (*conns, error) {
	db, err := connect(ctx, url)
	if err != nil {
		return nil, err
	}

	logger = logger.Set("pkg", "internal/database/psql")

	return &conns{
		db:     db,
		logger: logger,
	}, nil
}

func connect(ctx context.Context, url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("open DB conn: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping DB conn: %w", err)
	}

	return db, nil
}

func (c *conns) Conn(ctx context.Context) (Conn, error) {
	if !trx.Exists(ctx) {
		return adaptDB(c.db), nil
	}

	if t, ok := trx.GetValue(ctx, trxKey{}); ok {
		if sqlTx, ok := t.(*sql.Tx); ok {
			return adaptTx(sqlTx), nil
		} else {
			c.logger.Log(log.Warning).
				Set("message", "transaction is not *sql.Tx").
				Write()
		}
	}

	sqlTx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin sql tx: %w", err)
	}

	trx.SetValue(ctx, trxKey{}, sqlTx)
	trx.AddRollback(ctx, sqlTx.Rollback)
	trx.AddCommit(ctx, sqlTx.Commit)

	return adaptTx(sqlTx), nil
}

func (c *conns) Close() error {
	return c.db.Close()
}

// *sql.DB adapter

type dbAdapter struct {
	adaptee *sql.DB
}

func adaptDB(db *sql.DB) *dbAdapter {
	return &dbAdapter{adaptee: db}
}

func (a *dbAdapter) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	return a.adaptee.ExecContext(ctx, query, args...)
}

func (a *dbAdapter) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	return a.adaptee.QueryContext(ctx, query, args...)
}

func (a *dbAdapter) QueryRow(ctx context.Context, query string, args ...any) Row {
	return a.adaptee.QueryRowContext(ctx, query, args...)
}

// *sql.Tx adapter

type txAdapter struct {
	adaptee *sql.Tx
}

func adaptTx(tx *sql.Tx) *txAdapter {
	return &txAdapter{adaptee: tx}
}

func (a *txAdapter) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	return a.adaptee.ExecContext(ctx, query, args...)
}

func (a *txAdapter) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	return a.adaptee.QueryContext(ctx, query, args...)
}

func (a *txAdapter) QueryRow(ctx context.Context, query string, args ...any) Row {
	return a.adaptee.QueryRowContext(ctx, query, args...)
}
