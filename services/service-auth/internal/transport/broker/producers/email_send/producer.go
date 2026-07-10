package email_send

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/event"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/broker"
)

const topic = "email.send"

type client interface {
	Produce(ctx context.Context, msgs []broker.ProduceMessage) error
}

type Producer struct {
	client client
}

func NewProducer(client client) *Producer {
	return &Producer{
		client: client,
	}
}

func (p *Producer) Produce(ctx context.Context, events []event.EmailSend) error {
	msgs := make([]broker.ProduceMessage, 0, len(events))
	for _, e := range events {
		value, _ := json.Marshal(e)

		msgs = append(msgs, broker.ProduceMessage{
			Topic: topic,
			Key:   []byte(uuid.NewString()),
			Value: []byte(value),
		})
	}

	return p.client.Produce(ctx, msgs)
}
