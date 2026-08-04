package grpcutil

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
)

func ConvertError(err error, logger log.Logger) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, validate.ErrUnauthorized) {
		return Unauthenticated()
	}

	{
		var e *validate.Error
		if errors.As(err, &e) {
			return ValidationError(e)
		}
	}

	{
		var e validate.Errors
		if errors.As(err, &e) {
			return ValidationErrors(e)
		}
	}

	logger.Log(log.Error).
		Set("message", "unexpected error").
		Set("error", err.Error()).
		Write()

	return InternalError()
}

func InternalError() error {
	return status.Error(codes.Internal, "internal error")
}

func Unauthenticated() error {
	return status.Error(codes.Unauthenticated, "authorization error")
}
