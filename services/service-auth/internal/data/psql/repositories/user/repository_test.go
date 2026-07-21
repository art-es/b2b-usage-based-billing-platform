package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql/psqlmock"
)

func TestRepository_Find(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doScanRow := func(dest ...any) error {
			*(dest[0].(*string)) = "test_id"
			*(dest[1].(*string)) = "test_name"
			*(dest[2].(*string)) = "test_email"
			*(dest[3].(*string)) = "test_password_hash"
			*(dest[4].(*bool)) = true

			return nil
		}

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(5)...).
				DoAndReturn(doScanRow)

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(gomock.Any(), gomock.Eq(queryFind), gomock.Eq("test_id")).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr, err := repo.Find(ctx, "test_id")

		assert.NoError(t, err)
		assert.Equal(t, &user.User{
			ID:           "test_id",
			Name:         "test_name",
			Email:        "test_email",
			PasswordHash: "test_password_hash",
			IsVerified:   true,
		}, usr)
	})

	t.Run("not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doScanRow := func(...any) error {
			return sql.ErrNoRows
		}

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(5)...).
				DoAndReturn(doScanRow)

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(gomock.Any(), gomock.Eq(queryFind), gomock.Eq("test_id")).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr, err := repo.Find(ctx, "test_id")

		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, usr)
	})

	t.Run("query execute error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doScanRow := func(...any) error {
			return errors.New("test error")
		}

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(5)...).
				DoAndReturn(doScanRow)

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(gomock.Any(), gomock.Eq(queryFind), gomock.Eq("test_id")).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr, err := repo.Find(ctx, "test_id")

		assert.EqualError(t, err, "query execute: test error")
		assert.Nil(t, usr)
	})

	t.Run("get conn error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		mockConns.EXPECT().
			Conn(gomock.Any()).
			Return(nil, errors.New("test error"))

		repo := NewRepository(mockConns)
		usr, err := repo.Find(ctx, "test_id")

		assert.EqualError(t, err, "test error")
		assert.Nil(t, usr)
	})
}

func gomockAnyList(n int) []any {
	list := make([]any, n)
	for i := range n {
		list[i] = gomock.Any()
	}
	return list
}
