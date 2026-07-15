package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/send_email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/env"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql/repositories/email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/nats"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/retry"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/shutdown"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/broker/producers/email_send"
)

var (
	batchSize int
)

var (
	logger     log.Logger
	shutdowner *shutdown.Shutdowner
	usecase    *send_email_verification.Usecase
	repository *email_verification.Repository
)

func main() {
	flag.IntVar(&batchSize, "batch", 10, "batch size, default: 10")
	flag.Parse()

	logger = log.NewLogger(nil).Set("pkg", "cmd/mail_verification_cron")

	shutdowner = shutdown.New(logger)
	defer shutdowner.Shutdown()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := build(ctx); err != nil {
		logger.Log(log.Error).
			Set("message", "build error").
			Write()
		return
	}

	logger.Log(log.Info).
		Set("message", "cron started").
		Write()

	go func() {
		if err := run(ctx); err != nil {
			logger.Log(log.Error).
				Set("message", "run error").
				Set("error", err.Error()).
				Write()
		}
		stop()
	}()

	<-ctx.Done()

	logger.Log(log.Info).
		Set("message", "service finished").
		Write()
}

func build(ctx context.Context) error {
	if batchSize < 1 {
		return errors.New("-batch must be more than 0")
	}

	envs, err := env.ParseVars(
		env.Required(env.FieldPsqlUrl),
		env.Required(env.FieldNatsUrl),
	)
	if err != nil {
		return fmt.Errorf("parse env vars: %w", err)
	}

	psqlConn, err := psql.Connect(ctx, envs.Get(env.FieldPsqlUrl), logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	shutdowner.Add(psqlConn)

	natsConn, err := nats.Connect(envs.Get(env.FieldNatsUrl))
	if err != nil {
		return fmt.Errorf("connect nats: %w", err)
	}
	shutdowner.Add(natsConn)

	natsProducer := nats.NewProducer(natsConn)

	repository = email_verification.NewRepository(psqlConn)
	emailSendProducer := email_send.NewProducer(natsProducer)

	usecase = send_email_verification.NewUsecase(
		repository,
		emailSendProducer,
		logger,
		batchSize,
	)

	return nil
}

func run(ctx context.Context) error {
	err := repository.ClearDeprecated(ctx)
	if err != nil {
		logger.Log(log.Error).
			Set("message", "clear deprecated email verifications error").
			Set("error", err.Error()).
			Write()
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var updatesCount int

		err = retry.Retry(5, func() (err error) {
			updatesCount, err = usecase.Do(ctx)
			return
		})
		if err != nil {
			return err
		}

		if updatesCount == 0 {
			return nil
		}
	}
}
