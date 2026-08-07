package password_reset

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql"
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

	query := `INSERT INTO password_resets (user_id) VALUES ($1)`
	args := []any{userID}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) GetUnsent(ctx context.Context, batchSize int) ([]*user.PasswordReset, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT v.token, u.email
		FROM password_resets AS v
		JOIN users AS u ON u.id = v.user_id
		WHERE 
			v.sent_at IS NULL 
			AND u.verified_at IS NULL
		ORDER BY v.created_at
		LIMIT $1
		FOR UPDATE OF v SKIP LOCKED`
	args := []any{batchSize}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("execute query: %w", err)
	}

	defer rows.Close()

	resets := make([]*user.PasswordReset, 0, batchSize)

	for rows.Next() {
		reset := &user.PasswordReset{}
		if err := rows.Scan(&reset.Token, &reset.Email); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		resets = append(resets, reset)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration: %w", err)
	}

	return resets, nil
}

func (r *Repository) MarkAsSent(ctx context.Context, tokens []string) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE password_resets
		SET sent_at = current_timestamp
		WHERE token = ANY($1::uuid[])`
	args := []any{pq.Array(tokens)}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) ClearDeprecated(ctx context.Context) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM password_resets WHERE created_at + INTERVAL '3 days' > current_timestamp`

	_, err = conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) Find(ctx context.Context, token string) (*user.PasswordReset, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT r.token, u.user_id
		FROM password_resets AS r
		JOIN users AS u ON u.id = r.user_id
		WHERE r.token = $1 AND u.verified_at IS NULL
		FOR UPDATE OF r SKIP LOCKED`
	args := []any{token}

	var reset user.PasswordReset

	err = conn.QueryRow(ctx, query, args...).
		Scan(&reset.Token, &reset.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("query execute: %w", err)
	}

	return &reset, nil
}

func (r *Repository) DeleteForUser(ctx context.Context, userID string) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM password_resets WHERE user_id = $1`
	args := []any{userID}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}
