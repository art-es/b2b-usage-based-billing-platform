package login

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/hash"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx/trxutil"
)

type Usecase struct {
	jwtService        jwtService
	hashService       hashService
	timeService       timeService
	uuidService       uuidService
	sessionRepository sessionRepository
	userRepository    userRepository
	logger            log.Logger
	jwtSecret         []byte
}

func NewUsecase(
	jwtService jwtService,
	hashService hashService,
	timeService timeService,
	uuidService uuidService,
	sessionRepository sessionRepository,
	userRepository userRepository,
	logger log.Logger,
	jwtSecret []byte,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/login")

	return &Usecase{
		jwtService:        jwtService,
		hashService:       hashService,
		timeService:       timeService,
		uuidService:       uuidService,
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		logger:            logger,
		jwtSecret:         jwtSecret,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*dto.Response, error) {
	now := u.timeService.GetCurrentTime()

	usr, err := u.authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	return u.createSession(ctx, usr, now)
}

func (u *Usecase) authenticate(ctx context.Context, req *dto.Request) (*user.User, error) {
	usr, err := u.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, dto.ErrWrongCredentials
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	err = u.hashService.Compare(req.Password, usr.PasswordHash)
	if err != nil {
		if errors.Is(err, hash.ErrMismatch) {
			return nil, dto.ErrWrongCredentials
		}

		return nil, fmt.Errorf("compare password with hash: %w", err)
	}

	if !usr.IsVerified {
		return nil, dto.ErrEmailNotVerified
	}

	return usr, nil
}

func (u *Usecase) createSession(ctx context.Context, usr *user.User, now time.Time) (*dto.Response, error) {
	refreshToken := u.uuidService.Generate()
	refreshTokenHash, err := u.hashService.Generate(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token hash: %w", err)
	}

	var accessToken string

	ctx = trx.Begin(ctx)
	err = func() error {
		ses := session.NewSession(usr.ID, refreshTokenHash, now)
		err := u.sessionRepository.Save(ctx, ses)
		if err != nil {
			return fmt.Errorf("save session: %w", err)
		}

		accessToken, err = u.jwtService.Generate(u.jwtSecret, jwt.NewClaims(ses.ID, usr.ID))
		if err != nil {
			return fmt.Errorf("generate access token as jwt: %w", err)
		}

		return nil
	}()
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, err.Error())
		return nil, err
	}

	err = trx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit trx: %w", err)
	}

	return &dto.Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
