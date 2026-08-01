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

func (r *Repository) Get(ctx context.Context, userID string, cursor *session.ListCursor) ([]*session.Session, *session.ListCursor, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, nil, err
	}

	args := make([]any, 0, 3)
	args = append(args, userID)

	var whereCursor string
	if cursor != nil {
		whereCursor = "(created_at, id) < ($2, $3)"
		args = append(args, cursor.CreatedAt, cursor.ID)
	}

	const limit = 10
	query := fmt.Sprintf(
		`SELECT *
		FROM sessions
		WHERE user_id = $1%s
		ORDER BY created_at DESC, id DESC
		LIMIT %d`,
		whereCursor, limit+1,
	)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("query execute: %w", err)
	}

	defer rows.Close()

	var list []*session.Session

	for rows.Next() {
		s := &session.Session{}
		err = rows.Scan(
			&s.ID,
			&s.UserID,
			&s.OrganizationID,
			&s.RefreshTokenHash,
			&s.RefreshTokenExpiresAt,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("scan row: %w", err)
		}

		list = append(list, s)
	}

	hasMore := len(list) == limit+1
	list = list[:limit]
	nextCursor := session.GetNextCursor(list, hasMore)

	return list, nextCursor, nil
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
			refresh_token_expires_at,
			created_at
		FROM sessions
		WHERE refresh_token_hash = $1 
		FOR UPDATE SKIP LOCKED`
	args := []any{refreshTokenHash}

	s := &session.Session{}
	err = conn.QueryRow(ctx, query, args...).
		Scan(
			&s.ID,
			&s.UserID,
			&s.OrganizationID,
			&s.RefreshTokenHash,
			&s.RefreshTokenExpiresAt,
			&s.CreatedAt,
		)

	if err != nil {
		return nil, fmt.Errorf("query execute: %w", err)
	}

	return s, nil
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
