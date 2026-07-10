package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/kafka"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
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
	usecase    *email_verification.Usecase
)

func main() {
	flag.IntVar(&batchSize, "batch", 10, "batch size, default: 10")
	flag.Parse()

	logger = log.NewLogger(nil).Set("pkg", "cmd/mail_verification_cron")
	shutdowner = shutdown.NewManager(logger)
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

	psqlConn, err := psql.Connect(ctx, logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	shutdowner.Add(psqlConn)

	kafkaProducer, err := kafka.NewProducer(ctx)
	if err != nil {
		return fmt.Errorf("connect kafka: %w", err)
	}
	shutdowner.Add(kafkaProducer)

	// Repositories
	verificationRepository := repositories.NewVerificationRepository(psqlConn)

	// Broker
	emailSendProducer := email_send.NewProducer(kafkaProducer)

	usecase = email_verification.NewUsecase(
		verificationRepository,
		emailSendProducer,
		logger,
		batchSize,
	)

	return nil
}

func run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var updates int

		err := retry.Retry(5, func() (err error) {
			updates, err = usecase.Do(ctx)
			return
		})
		if err != nil {
			return err
		}

		if updates == 0 {
			return nil
		}
	}
}
