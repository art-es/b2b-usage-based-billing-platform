package nats

import (
	"errors"
	"fmt"
	"os"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/broker"
	"github.com/nats-io/nats.go"
)

type Conn struct {
	conn *nats.Conn
	js   nats.JetStream
}

func Connect() (*Conn, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return nil, errors.New("NATS_URL required")
	}

	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to nats server: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create jetstream: %w", err)
	}

	for _, subject := range broker.SupportedSubjects() {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     subject,
			Subjects: []string{subject},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			return nil, fmt.Errorf("add stream %q: %w", subject, err)
		}
	}

	return &Conn{
		conn: conn,
		js:   js,
	}, nil
}

func (c *Conn) Close() error {
	c.conn.Close()
	return nil
}
