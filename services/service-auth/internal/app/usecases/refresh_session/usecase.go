package refresh_session

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/refresh_session/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx/trxutil"
)

type Usecase struct {
	jwtService             jwtService
	keyedHashService       keyedHashService
	uuidService            uuidService
	sessionRepository      sessionRepository
	jwtSecret              []byte
	refreshTokenHashSecret []byte
	logger                 log.Logger
}

func NewUsecase(
	jwtService jwtService,
	keyedHashService keyedHashService,
	uuidService uuidService,
	sessionRepository sessionRepository,
	jwtSecret string,
	refreshTokenHashSecret string,
	logger log.Logger,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/refresh_session")

	return &Usecase{
		jwtService:             jwtService,
		keyedHashService:       keyedHashService,
		uuidService:            uuidService,
		sessionRepository:      sessionRepository,
		jwtSecret:              []byte(jwtSecret),
		refreshTokenHashSecret: []byte(refreshTokenHashSecret),
		logger:                 logger,
	}
}

func (u *Usecase) Do(ctx context.Context, refreshToken string) (*dto.Response, error) {
	refreshTokenHash, err := u.keyedHashService.Generate(u.refreshTokenHashSecret, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("generate input refresh token hash: %w", err)
	}

	res := &dto.Response{}
	ctx = trx.Begin(ctx)
	err = func() error {
		ses, err := u.sessionRepository.GetByRefreshTokenHash(ctx, refreshTokenHash)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return dto.ErrInvalidToken
			}

			return fmt.Errorf("get by refresh token: %w", err)
		}

		res.RefreshToken = u.uuidService.Generate()
		ses.RefreshTokenHash, err = u.keyedHashService.Generate(u.refreshTokenHashSecret, res.RefreshToken)
		if err != nil {
			return fmt.Errorf("generate refresh token hash: %w", err)
		}

		err = u.sessionRepository.Save(ctx, ses)
		if err != nil {
			return fmt.Errorf("save session: %w", err)
		}

		res.AccessToken, err = u.jwtService.Generate(u.jwtSecret, jwt.NewClaims(ses.ID, ses.UserID))
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

	return res, nil
}
