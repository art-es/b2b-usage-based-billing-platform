package get_me

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_me/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()

	const (
		testSessionID = "test_session_id"
		testUserID    = "test_user_id"
		testUserEmail = "test_user_email"
		testUserName  = "test_user_name"
		testOrgnID    = "test_orgn_id"
		testOrgnName  = "test_orgn_name"
	)

	type deps struct {
		mockUserRepository *MockuserRepository
		mockOrgnRepository *MockorgnRepository
		usecase            *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockUserRepository := NewMockuserRepository(mockCtrl)
		mockOrgnRepository := NewMockorgnRepository(mockCtrl)

		return &deps{
			mockUserRepository: mockUserRepository,
			mockOrgnRepository: mockOrgnRepository,
			usecase:            NewUsecase(mockUserRepository, mockOrgnRepository),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		expUser := &user.User{
			ID:    testUserID,
			Name:  testUserName,
			Email: testUserEmail,
		}

		d.mockUserRepository.EXPECT().
			Find(gomock.Any(), gomock.Eq(testUserID)).
			Return(expUser, nil)

		expOrgn := &orgn.Orgn{
			ID:   testOrgnID,
			Name: testOrgnName,
		}

		d.mockOrgnRepository.EXPECT().
			Find(gomock.Any(), gomock.Eq(testOrgnID)).
			Return(expOrgn, nil)

		res, err := d.usecase.Do(ctx, &jwt.Claims{
			SessionID: testSessionID,
			UserID:    testUserID,
			OrgnID:    ptr.To(string(testOrgnID)),
		})

		assert.NoError(t, err)

		expRes := &dto.Response{
			SessionID: testSessionID,
			User: dto.ResponseUser{
				Email: testUserEmail,
				Name:  testUserName,
			},
			Orgn: &dto.ResponseOrgn{
				ID:   testOrgnID,
				Name: testOrgnName,
			},
		}

		assert.Equal(t, expRes, res)
	})

	t.Run("ok, no orgn id", func(t *testing.T) {
		d := newDeps()

		expUser := &user.User{
			ID:    testUserID,
			Name:  testUserName,
			Email: testUserEmail,
		}

		d.mockUserRepository.EXPECT().
			Find(gomock.Any(), gomock.Eq(testUserID)).
			Return(expUser, nil)

		res, err := d.usecase.Do(ctx, &jwt.Claims{
			SessionID: testSessionID,
			UserID:    testUserID,
			OrgnID:    nil,
		})

		assert.NoError(t, err)

		expRes := &dto.Response{
			SessionID: testSessionID,
			User: dto.ResponseUser{
				Email: testUserEmail,
				Name:  testUserName,
			},
			Orgn: nil,
		}

		assert.Equal(t, expRes, res)
	})

	t.Run("find orgn error", func(t *testing.T) {
		d := newDeps()

		expUser := &user.User{
			ID:    testUserID,
			Name:  testUserName,
			Email: testUserEmail,
		}

		d.mockUserRepository.EXPECT().
			Find(gomock.Any(), gomock.Eq(testUserID)).
			Return(expUser, nil)

		d.mockOrgnRepository.EXPECT().
			Find(gomock.Any(), gomock.Eq(testOrgnID)).
			Return(nil, errors.New("test error"))

		res, err := d.usecase.Do(ctx, &jwt.Claims{
			SessionID: testSessionID,
			UserID:    testUserID,
			OrgnID:    ptr.To(string(testOrgnID)),
		})

		assert.EqualError(t, err, "find orgn: test error")
		assert.Nil(t, res)
	})

	t.Run("find user error", func(t *testing.T) {
		d := newDeps()

		d.mockUserRepository.EXPECT().
			Find(gomock.Any(), gomock.Eq(testUserID)).
			Return(nil, errors.New("test error"))

		res, err := d.usecase.Do(ctx, &jwt.Claims{
			SessionID: testSessionID,
			UserID:    testUserID,
			OrgnID:    ptr.To(string(testOrgnID)),
		})

		assert.EqualError(t, err, "find user: test error")
		assert.Nil(t, res)
	})
}
