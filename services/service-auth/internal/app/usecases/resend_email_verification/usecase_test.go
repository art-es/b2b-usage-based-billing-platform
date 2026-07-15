package resend_email_verification

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/resend_email_verification/dto"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()

	type deps struct {
		mockEmailVerificationRepository *MockemailVerificationRepository
		usecase                         *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockEmailVerificationRepository := NewMockemailVerificationRepository(mockCtrl)

		return &deps{
			mockEmailVerificationRepository: mockEmailVerificationRepository,
			usecase:                         NewUsecase(mockEmailVerificationRepository),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationRepository.EXPECT().
			HasUnsentByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(false, nil)

		d.mockEmailVerificationRepository.EXPECT().
			CreateForEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(nil)

		err := d.usecase.Do(ctx, "test@example.com")
		assert.NoError(t, err)
	})

	t.Run("has unsent verification", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationRepository.EXPECT().
			HasUnsentByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(true, nil)

		err := d.usecase.Do(ctx, "test@example.com")
		assert.NoError(t, err)
	})

	t.Run("email is already verified", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationRepository.EXPECT().
			HasUnsentByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(false, fmt.Errorf("repo: %w", user.ErrEmailVerified))

		err := d.usecase.Do(ctx, "test@example.com")
		assert.EqualError(t, err, "email is already verified")
		assert.ErrorIs(t, err, dto.ErrEmailVerified)
	})

	t.Run("invalid email", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationRepository.EXPECT().
			HasUnsentByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(false, fmt.Errorf("repo: %w", user.ErrUserNotFound))

		err := d.usecase.Do(ctx, "test@example.com")
		assert.EqualError(t, err, "invalid email")
		assert.ErrorIs(t, err, dto.ErrInvalidEmail)
	})

	t.Run("check unsent verification error", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationRepository.EXPECT().
			HasUnsentByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(false, errors.New("test error"))

		err := d.usecase.Do(ctx, "test@example.com")
		assert.EqualError(t, err, "check existence of unsent verification: test error")
	})

	t.Run("create verification error", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationRepository.EXPECT().
			HasUnsentByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(false, nil)

		d.mockEmailVerificationRepository.EXPECT().
			CreateForEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(errors.New("test error"))

		err := d.usecase.Do(ctx, "test@example.com")
		assert.EqualError(t, err, "create verification for email: test error")
	})
}
