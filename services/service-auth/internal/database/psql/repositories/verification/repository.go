package verification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql"
	"github.com/lib/pq"
)

type Repository struct {
	conns psql.Conns
}

func NewRepository(conns psql.Conns) *Repository {
	return &Repository{
		conns: conns,
	}
}

func (r *Repository) Create(ctx context.Context, userID string) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO verifications (user_id) VALUES ($1)`
	args := []any{userID}

	_, err = conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) GetUnsent(ctx context.Context, batchSize int) ([]*user.Verification, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT v.token, u.email
		FROM verifications AS v
		JOIN users AS u ON u.id = v.user_id
		WHERE v.email_sent_at is NULL
		ORDER BY v.created_at
		LIMIT $1
		FOR UPDATE OF v SKIP LOCKED`
	args := []any{batchSize}

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("execute query: %w", err)
	}

	defer rows.Close()

	vers := make([]*user.Verification, 0, batchSize)

	for rows.Next() {
		ver := &user.Verification{}
		if err := rows.Scan(ver.Token, ver.Email); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		vers = append(vers, ver)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration: %w", err)
	}

	return vers, nil
}

func (r *Repository) MarkAsSent(ctx context.Context, tokens []string) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE verifications
		SET email_sent_at = current_timestamp
		WHERE token = ANY($1)`
	args := []any{pq.Array(tokens)}

	_, err = conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}
