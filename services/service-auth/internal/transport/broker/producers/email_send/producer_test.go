//go:generate mockgen -source=usecase.go -destination=usecase_mock_test.go -package=$GOPACKAGE
package email_send

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/event"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/broker"
)

func TestProducer(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockClient := NewMockclient(mockCtrl)

		expMsg := broker.ProduceMessage{
			Subject:        broker.SubjectEmailSend,
			IdempotencyKey: "test-idempotency-key",
			Payload: []byte(`{
				"email": "test-email",
				"subject": "Test Subject",
				"content": "Test Content"
			}`),
		}

		mockClient.EXPECT().
			Produce(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, actMsg broker.ProduceMessage) {
				assert.Equal(t, expMsg.Subject, actMsg.Subject)
				assert.Equal(t, expMsg.IdempotencyKey, actMsg.IdempotencyKey)
				assert.JSONEq(t, string(expMsg.Payload), string(actMsg.Payload))
			}).
			Return(nil)

		producer := NewProducer(mockClient)

		err := producer.Produce(ctx, event.EmailSend{
			IdempotencyKey: "test-idempotency-key",
			Email:          "test-email",
			Subject:        "Test Subject",
			Content:        "Test Content",
		})

		assert.NoError(t, err)
	})

	t.Run("client error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockClient := NewMockclient(mockCtrl)

		expMsg := broker.ProduceMessage{
			Subject:        broker.SubjectEmailSend,
			IdempotencyKey: "test-idempotency-key",
			Payload: []byte(`{
				"email": "test-email",
				"subject": "Test Subject",
				"content": "Test Content"
			}`),
		}

		mockClient.EXPECT().
			Produce(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, actMsg broker.ProduceMessage) {
				assert.Equal(t, expMsg.Subject, actMsg.Subject)
				assert.Equal(t, expMsg.IdempotencyKey, actMsg.IdempotencyKey)
				assert.JSONEq(t, string(expMsg.Payload), string(actMsg.Payload))
			}).
			Return(errors.New("test error"))

		producer := NewProducer(mockClient)

		err := producer.Produce(ctx, event.EmailSend{
			IdempotencyKey: "test-idempotency-key",
			Email:          "test-email",
			Subject:        "Test Subject",
			Content:        "Test Content",
		})

		assert.EqualError(t, err, "client: test error")
	})
}
