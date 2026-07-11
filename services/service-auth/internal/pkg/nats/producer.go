package nats

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/broker"
)

type Producer struct {
	conn *Conn
}

func NewProducer(conn *Conn) *Producer {
	return &Producer{conn: conn}
}

func (p *Producer) Produce(ctx context.Context, msg broker.ProduceMessage) error {
	_, err := p.conn.js.Publish(
		msg.Subject,
		msg.Payload,
		nats.Context(ctx),
		nats.MsgId(msg.IdempotencyKey),
	)

	return err
}
