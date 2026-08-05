package sessions

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql"
)

type Repository struct {
	conns psql.Conns
}

func NewRepository(conns psql.Conns) *Repository {
	return &Repository{
		conns: conns,
	}
}

func (r *Repository) Save(ctx context.Context, ses *session.Session) error {
	if ses.Stored() {
		return r.update(ctx, ses)
	}

	return r.insert(ctx, ses)
}

func (r *Repository) insert(ctx context.Context, s *session.Session) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO sessions (
			user_id, 
			refresh_token_hash, 
			refresh_token_expires_at
		) 
		VALUES ($1, $2, $3)
		RETURNING id`
	args := []any{s.UserID, s.RefreshTokenHash, s.RefreshTokenExpiresAt}

	err = conn.QueryRow(ctx, query, args...).
		Scan(&s.ID)

	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) update(ctx context.Context, s *session.Session) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE sessions
		SET 
			refresh_token_hash = $1,
			refresh_token_expires_at = $2,
			updated_at = current_timestamp
		WHERE id = $3`
	args := []any{s.RefreshTokenHash, s.RefreshTokenExpiresAt, s.ID}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) DeleteAll(ctx context.Context, userID string) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM sessions WHERE user_id = $1`
	args := []any{userID}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, userID, sessionID string) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM sessions WHERE user_id = $1 AND session_id = $2`
	args := []any{userID, sessionID}

	res, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if n == 0 {
		return repository.ErrNotFound
	}

	return nil
}
