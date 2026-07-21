package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql/psqlmock"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql/psqlutil"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
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

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(5)...).
				Return(sql.ErrNoRows)

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

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(5)...).
				Return(errors.New("test error"))

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

func TestRepository_FindByEmail(t *testing.T) {
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
				QueryRow(gomock.Any(), gomock.Eq(queryFindByEmail), gomock.Eq("test_email")).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr, err := repo.FindByEmail(ctx, "test_email")

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

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(5)...).
				Return(sql.ErrNoRows)

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(gomock.Any(), gomock.Eq(queryFindByEmail), gomock.Eq("test_email")).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr, err := repo.FindByEmail(ctx, "test_email")

		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, usr)
	})

	t.Run("query execute error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(5)...).
				Return(errors.New("test error"))

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(gomock.Any(), gomock.Eq(queryFindByEmail), gomock.Eq("test_email")).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr, err := repo.FindByEmail(ctx, "test_email")

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
		usr, err := repo.FindByEmail(ctx, "test_email")

		assert.EqualError(t, err, "test error")
		assert.Nil(t, usr)
	})
}

func TestRepository_Save(t *testing.T) {
	ctx := context.Background()
	testUnstoredUser := user.User{
		Name:         "test_name",
		Email:        "test_email",
		PasswordHash: "test_password_hash",
	}

	t.Run("insert: ok", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doScanRow := func(dest ...any) error {
			*(dest[0].(*string)) = "test_id"

			return nil
		}

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(1)...).
				DoAndReturn(doScanRow)

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(
					gomock.Any(),
					gomock.Eq(queryInsert),
					gomock.Eq("test_name"),
					gomock.Eq("test_email"),
					gomock.Eq("test_password_hash"),
				).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr := ptr.To(testUnstoredUser)
		err := repo.Save(ctx, usr)

		assert.NoError(t, err)
		assert.Equal(t, &user.User{
			ID:           "test_id",
			Name:         "test_name",
			Email:        "test_email",
			PasswordHash: "test_password_hash",
		}, usr)
	})

	t.Run("insert: unique error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(1)...).
				Return(&pq.Error{Code: psqlutil.UniqueViolationErrorCode})

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(
					gomock.Any(),
					gomock.Eq(queryInsert),
					gomock.Eq("test_name"),
					gomock.Eq("test_email"),
					gomock.Eq("test_password_hash"),
				).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr := ptr.To(testUnstoredUser)
		err := repo.Save(ctx, usr)

		assert.ErrorIs(t, err, repository.ErrUnique)
		assert.Equal(t, testUnstoredUser, *usr) // assert no changes
	})

	t.Run("insert: query execute error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doQueryRow := func(context.Context, string, ...any) psql.Row {
			mockRow := psqlmock.NewMockRow(mockCtrl)
			mockRow.EXPECT().
				Scan(gomockAnyList(1)...).
				Return(errors.New("test error"))

			return mockRow
		}

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				QueryRow(
					gomock.Any(),
					gomock.Eq(queryInsert),
					gomock.Eq("test_name"),
					gomock.Eq("test_email"),
					gomock.Eq("test_password_hash"),
				).
				DoAndReturn(doQueryRow)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		usr := ptr.To(testUnstoredUser)
		err := repo.Save(ctx, usr)

		assert.EqualError(t, err, "query execute: test error")
		assert.Equal(t, testUnstoredUser, *usr) // assert no changes
	})

	t.Run("insert: get conn error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		mockConns.EXPECT().
			Conn(gomock.Any()).
			Return(nil, errors.New("test error"))

		repo := NewRepository(mockConns)
		usr := ptr.To(testUnstoredUser)
		err := repo.Save(ctx, usr)

		assert.EqualError(t, err, "test error")
		assert.Equal(t, testUnstoredUser, *usr) // assert no changes
	})

	t.Run("update: not implemented", func(t *testing.T) {
		repo := NewRepository(nil)

		usr := ptr.To(testUnstoredUser)
		usr.ID = "test_id"
		assert.True(t, usr.Stored())

		err := repo.Save(ctx, usr)
		assert.EqualError(t, err, "UPDATE not implemented")
	})
}

func TestRepository_MarkAsVerified(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				Exec(
					gomock.Any(),
					gomock.Eq(queryMarkAsVerified),
					gomock.Eq("test_id"),
				).
				Return(nil, nil)

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		err := repo.MarkAsVerified(ctx, "test_id")
		assert.NoError(t, err)
	})

	t.Run("query execute error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		doGetConn := func(context.Context) (psql.Conn, error) {
			mockConn := psqlmock.NewMockConn(mockCtrl)
			mockConn.EXPECT().
				Exec(
					gomock.Any(),
					gomock.Eq(queryMarkAsVerified),
					gomock.Eq("test_id"),
				).
				Return(nil, errors.New("test error"))

			return mockConn, nil
		}

		mockConns.EXPECT().
			Conn(gomock.Any()).
			DoAndReturn(doGetConn)

		repo := NewRepository(mockConns)
		err := repo.MarkAsVerified(ctx, "test_id")
		assert.EqualError(t, err, "query execute: test error")
	})

	t.Run("get conn error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockConns := psqlmock.NewMockConns(mockCtrl)

		mockConns.EXPECT().
			Conn(gomock.Any()).
			Return(nil, errors.New("test error"))

		repo := NewRepository(mockConns)
		err := repo.MarkAsVerified(ctx, "test_id")
		assert.EqualError(t, err, "test error")
	})
}

func gomockAnyList(n int) []any {
	list := make([]any, n)
	for i := range n {
		list[i] = gomock.Any()
	}
	return list
}
