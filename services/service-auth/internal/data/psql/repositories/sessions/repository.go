package sessions

import (
	"context"
	"errors"
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

func (r *Repository) Save(ctx context.Context, ses *session.Session) error {
	if ses.Stored() {
		return errors.New("UPDATE not implemented")
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
