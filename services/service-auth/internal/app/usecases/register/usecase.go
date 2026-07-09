//go:generate mockgen -source=usecase.go -destination=usecase_mock_test.go -package=$GOPACKAGE
package register

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
)

type userRepository interface {
	Store(ctx context.Context, usr *user.User) error
}

type hashService interface {
	Generate(s string) (string, error)
}

type Usecase struct {
	userRepository userRepository
	hashService    hashService
}

func NewUsecase(userRepository userRepository, hashService hashService) *Usecase {
	return &Usecase{
		userRepository: userRepository,
		hashService:    hashService,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) error {
	passwordHash, err := u.hashService.Generate(req.Password)
	if err != nil {
		return fmt.Errorf("generate password hash: %w", err)
	}

	usr := user.NewRegisteredUser(req.Name, req.Email, passwordHash)

	err = u.userRepository.Store(ctx, usr)
	if err != nil {
		if errors.Is(err, repository.ErrUnique) {
			return dto.ErrEmailInUse
		}

		return fmt.Errorf("store user: %w", err)
	}

	return nil
}
