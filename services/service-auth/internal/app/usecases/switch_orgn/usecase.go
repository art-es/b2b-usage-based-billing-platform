package switch_orgn

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/trx/trxutil"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/switch_orgn/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
)

type sessionRepository interface {
	FindByID(ctx context.Context, id string) (*session.Session, error)
	Save(ctx context.Context, ses *session.Session) error
}

type orgnService interface {
	GetByID(ctx context.Context, orgnID string) (*orgn.Orgn, error)
}

type jwtService interface {
	Generate(secret []byte, claims *jwt.Claims) (string, error)
}

type Usecase struct {
	sessionRepository sessionRepository
	orgnService       orgnService
	jwtService        jwtService
	jwtSecret         []byte
	logger            log.Logger
}

func NewUsecase(
	sessionRepository sessionRepository,
	orgnService orgnService,
	jwtService jwtService,
	jwtSecret string,
	logger log.Logger,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/switch_orgn")

	return &Usecase{
		sessionRepository: sessionRepository,
		orgnService:       orgnService,
		jwtService:        jwtService,
		jwtSecret:         []byte(jwtSecret),
		logger:            logger,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*dto.Response, error) {
	org, err := u.orgnService.GetByID(ctx, req.OrgnID)
	if err != nil {
		if errors.Is(err, orgn.ErrNotFound) {
			return nil, errIncorrectOrgnID
		}

		return nil, fmt.Errorf("get orgn by id: %w", err)
	}

	if org.UserID != req.Auth.UserID {
		return nil, errIncorrectOrgnID
	}

	var ses *session.Session
	ctx = trx.Begin(ctx)
	err = func() error {
		ses, err = u.sessionRepository.FindByID(ctx, req.Auth.SessionID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return validate.ErrUnauthorized
			}

			return fmt.Errorf("find session for update: %w", err)
		}

		ses.SetOrgn(req.OrgnID)

		err = u.sessionRepository.Save(ctx, ses)
		if err != nil {
			return fmt.Errorf("save session: %w", err)
		}

		return nil
	}()
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, err.Error())
		return nil, err
	}

	accessToken, err := u.jwtService.Generate(u.jwtSecret, jwt.NewClaims(ses.ID, ses.UserID, ses.OrganizationID))
	if err != nil {
		return nil, fmt.Errorf("generate access token as jwt: %w", err)
	}

	return &dto.Response{
		AccessToken: accessToken,
	}, nil
}
