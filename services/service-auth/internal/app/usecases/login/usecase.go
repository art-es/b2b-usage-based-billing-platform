package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/services/hash"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login/dto"
)

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

type sessionRepository interface {
	Create(ctx context.Context, sess *session.Session) error
}

type hashService interface {
	Compare(s string, hash string) error
	Generate(s string) (string, error)
}

type jwtService interface {
	Generate(sess *session.Session) (string, error)
}

type Usecase struct {
	userRepository    userRepository
	sessionRepository sessionRepository
	hashService       hashService
	jwtService        jwtService
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*session.Tokens, error) {
	usr, err := u.authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	return u.createSession(ctx, usr)
}

func (u *Usecase) authenticate(ctx context.Context, req *dto.Request) (*user.User, error) {
	usr, err := u.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, dto.ErrWrongCredentials
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if !usr.IsVerified {
		return nil, dto.ErrEmailNotVerified
	}

	err = u.hashService.Compare(req.Password, usr.PasswordHash)
	if err != nil {
		if errors.Is(err, hash.ErrMismatch) {
			return nil, dto.ErrWrongCredentials
		}

		return nil, fmt.Errorf("compare password with hash: %w", err)
	}

	return usr, nil
}

func (u *Usecase) createSession(ctx context.Context, usr *user.User) (*session.Tokens, error) {
	sess := session.NewSession(usr.ID)

	accessToken, err := u.jwtService.Generate(sess)
	if err != nil {
		return nil, fmt.Errorf("generate access token as JWT: %w", err)
	}

	refreshToken := session.NewRefreshToken()
	refreshTokenHash, err := u.hashService.Generate(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token hash: %w", err)
	}

	sess.RefreshTokenHash = refreshTokenHash
	err = u.sessionRepository.Create(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &session.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
