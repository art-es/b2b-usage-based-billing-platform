package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql/psqlutil"
)

const (
	queryFind = `
		SELECT id, name, email, password_hash, verified_at IS NOT NULL AS is_verified
		FROM users WHERE id = $1`

	queryFindByEmail = `
		SELECT id, name, email, password_hash, verified_at IS NOT NULL AS is_verified
		FROM users WHERE email = $1`

	queryInsert = `
		INSERT INTO users (name, email, password_hash) 
		VALUES ($1, $2, $3)`

	queryMarkAsVerified = `
		UPDATE users 
		SET verified_at = current_timestamp 
		WHERE id = $1`
)

type Repository struct {
	conns psql.Conns
}

func NewRepository(conns psql.Conns) *Repository {
	return &Repository{
		conns: conns,
	}
}

func (r *Repository) Find(ctx context.Context, id string) (*user.User, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	usr := &user.User{}
	err = conn.QueryRow(ctx, queryFind, id).
		Scan(&usr.ID, &usr.Name, &usr.Email, &usr.PasswordHash, &usr.IsVerified)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("query execute: %w", err)
	}

	return usr, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return nil, err
	}

	usr := &user.User{}
	err = conn.QueryRow(ctx, queryFindByEmail, email).
		Scan(&usr.ID, &usr.Name, &usr.Email, &usr.PasswordHash, &usr.IsVerified)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("query execute: %w", err)
	}

	return usr, nil
}

func (r *Repository) Save(ctx context.Context, usr *user.User) error {
	if usr.Stored() {
		return errors.New("UPDATE not implemented")
	}

	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	err = conn.QueryRow(ctx, queryInsert, usr.Name, usr.Email, usr.PasswordHash).
		Scan(&usr.ID)

	if err != nil {
		if psqlutil.IsUniqueViolationError(err) {
			return repository.ErrUnique
		}

		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}

func (r *Repository) MarkAsVerified(ctx context.Context, userID string) error {
	conn, err := r.conns.Conn(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, queryMarkAsVerified, userID)
	if err != nil {
		return fmt.Errorf("query execute: %w", err)
	}

	return nil
}
