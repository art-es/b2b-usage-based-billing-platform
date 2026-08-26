package orgn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/data/psql"
)

const columns = `id, name, user_id, created_at`

func scanColumns(o *orgn.Orgn) []any {
	return []any{&o.ID, &o.Name, &o.UserID, &o.CreatedAt}
}

type Repository struct {
	conns psql.Conns
}

func NewRepository(conns psql.Conns) *Repository {
	return &Repository{conns: conns}
}

func (r *Repository) Find(ctx context.Context, orgnID, userID string) (*orgn.Orgn, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT %s FROM organizations WHERE id = $1 AND user_id = $2`, columns)
	args := []any{orgnID, userID}

	if conn.IsTx() {
		query += " FOR UPDATE"
	}

	org := &orgn.Orgn{}
	err = conn.QueryRow(ctx, query, args...).Scan(scanColumns(org)...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("query execute: %w", err)
	}

	return org, nil
}

func (r *Repository) Get(ctx context.Context, userID string, cursor *orgn.ListCursor) ([]*orgn.Orgn, *orgn.ListCursor, error) {
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

	query := fmt.Sprintf(
		`SELECT %s
		FROM organizations
		WHERE user_id = $1%s
		ORDER BY created_at DESC, id DESC
		LIMIT %d`,
		columns, whereCursor, orgn.DBListLimit,
	)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}

		return nil, nil, fmt.Errorf("query execute: %w", err)
	}

	defer rows.Close()

	var list []*orgn.Orgn

	for rows.Next() {
		s := &orgn.Orgn{}
		err = rows.Scan(scanColumns(s)...)
		if err != nil {
			return nil, nil, fmt.Errorf("scan row: %w", err)
		}
		list = append(list, s)
	}

	list, nextCursor := orgn.HandleList(list)
	return list, nextCursor, nil
}

func (r *Repository) Save(ctx context.Context, o *orgn.Orgn) error {
	if o.Stored() {
		return r.update(ctx, o)
	}
	return r.insert(ctx, o)
}

func (r *Repository) insert(ctx context.Context, o *orgn.Orgn) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO orgns (name, user_id) VALUES ($1, $2) RETURNING id`
	args := []any{o.Name, o.UserID}

	err = conn.QueryRow(ctx, query, args...).Scan(&o.ID)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (r *Repository) update(ctx context.Context, o *orgn.Orgn) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE orgns SET name = $2 WHERE id = $1`
	args := []any{o.UserID, o.Name}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}
