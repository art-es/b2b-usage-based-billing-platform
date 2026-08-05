package sessions

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
)

const columns = `id, user_id, organization_id, refresh_token_hash, refresh_token_expires_at, created_at`

func scanColumns(s *session.Session) []any {
	return []any{&s.ID, &s.UserID, &s.OrganizationID, &s.RefreshTokenHash, &s.RefreshTokenExpiresAt, &s.CreatedAt}
}

func (r *Repository) find(ctx context.Context, query string, args []any) (*session.Session, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	s := &session.Session{}
	err = conn.QueryRow(ctx, query, args...).Scan(scanColumns(s)...)

	if err != nil {
		return nil, fmt.Errorf("query execute: %w", err)
	}

	return s, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*session.Session, error) {
	query := "SELECT " + columns + " FROM sessions WHERE id = $1 FOR UPDATE SKIP LOCKED"
	args := []any{id}

	return r.find(ctx, query, args)
}

func (r *Repository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*session.Session, error) {
	query := "SELECT " + columns + " FROM sessions WHERE refresh_token_hash = $1 FOR UPDATE SKIP LOCKED"
	args := []any{refreshTokenHash}

	return r.find(ctx, query, args)
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
		`SELECT %s
		FROM sessions
		WHERE user_id = $1%s
		ORDER BY created_at DESC, id DESC
		LIMIT %d`,
		columns, whereCursor, limit+1,
	)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("query execute: %w", err)
	}

	defer rows.Close()

	var list []*session.Session

	for rows.Next() {
		s := &session.Session{}
		err = rows.Scan(scanColumns(s)...)
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
