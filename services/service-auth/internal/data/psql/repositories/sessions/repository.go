package sessions

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
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

func (r *Repository) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*session.Session, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			id, 
			user_id, 
			organization_id, 
			refresh_token_hash,
			refresh_token_expires_at
		FROM sessions
		WHERE refresh_token_hash = $1`
	args := []any{refreshTokenHash}

	ses := &session.Session{}
	err = conn.QueryRowContext(ctx, query, args...).
		Scan(
			&ses.ID,
			&ses.UserID,
			&ses.OrganizationID,
			&ses.RefreshTokenHash,
			&ses.RefreshTokenExpiresAt,
		)
	if err != nil {
		return nil, fmt.Errorf("query execute: %w", err)
	}

	return ses, nil
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

	err = conn.QueryRowContext(ctx, query, args...).Scan(&s.ID)
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

	_, err = conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}
