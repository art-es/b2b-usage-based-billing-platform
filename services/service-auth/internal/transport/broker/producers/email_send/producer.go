package email_send

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/event"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/broker"
)

type client interface {
	Produce(ctx context.Context, msg broker.ProduceMessage) error
}

type messagePayload struct {
	Email   string `json:"email,omitempty"`
	Subject string `json:"subject,omitempty"`
	Content string `json:"content,omitempty"`
}

type Producer struct {
	client client
}

func NewProducer(client client) *Producer {
	return &Producer{
		client: client,
	}
}

func (p *Producer) Produce(ctx context.Context, ev event.EmailSend) error {
	payload, err := json.Marshal(messagePayload{
		Email:   ev.Email,
		Subject: ev.Subject,
		Content: ev.Content,
	})
	if err != nil {
		return fmt.Errorf("encode payload to json: %w", err)
	}

	return p.client.Produce(ctx, broker.ProduceMessage{
		Subject:        broker.SubjectEmailSend,
		IdempotencyKey: ev.IdempotencyKey,
		Payload:        payload,
	})
}
