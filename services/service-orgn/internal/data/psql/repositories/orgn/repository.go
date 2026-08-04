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

type Repository struct {
	conns psql.Conns
}

func NewRepository(conns psql.Conns) *Repository {
	return &Repository{conns: conns}
}

func (r *Repository) Find(ctx context.Context, id string) (*orgn.Orgn, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, user_id FROM organizations WHERE id = $1`
	args := []any{id}

	org := &orgn.Orgn{}
	err = conn.QueryRow(ctx, query, args...).
		Scan(
			&org.ID,
			&org.Name,
			&org.UserID,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("query execute: %w", err)
	}

	return org, nil
}
