package register

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()

	type deps struct {
		mockUserRepository *MockuserRepository
		mockHashService    *MockhashService
		usecase            *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockUserRepository := NewMockuserRepository(mockCtrl)
		mockHashService := NewMockhashService(mockCtrl)

		return &deps{
			mockUserRepository: mockUserRepository,
			mockHashService:    mockHashService,
			usecase:            NewUsecase(mockUserRepository, mockHashService),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Store(gomock.Any(), gomock.Eq(expUser)).
			Return(nil)

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("store user error", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Store(gomock.Any(), gomock.Eq(expUser)).
			Return(errors.New("test error"))

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "store user: test error")
	})

	t.Run("email is already in use", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Store(gomock.Any(), gomock.Eq(expUser)).
			Return(fmt.Errorf("repo: %w", repository.ErrUnique))

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "email is already in use")
		assert.ErrorIs(t, err, dto.ErrEmailInUse)
	})

	t.Run("generate password hash error", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("", errors.New("test error"))

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "generate password hash: test error")
	})
}
