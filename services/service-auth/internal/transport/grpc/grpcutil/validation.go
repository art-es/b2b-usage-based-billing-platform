package grpcutil

import (
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
)

const (
	validationErrorReasons = 1000
)

func ValidationErrors(vErrs validate.Errors) error {
	details := &errdetails.BadRequest{}
	for _, vErr := range vErrs {
		details.FieldViolations = append(details.FieldViolations, toValidatonError(vErr))
	}

	return newInvalidArgumentError(details)
}

func ValidationError(vErr *validate.Error) error {
	if vErr == nil {
		return nil
	}

	details := &errdetails.BadRequest{}
	details.FieldViolations = append(details.FieldViolations, toValidatonError(vErr))

	return newInvalidArgumentError(details)
}

func newInvalidArgumentError(details *errdetails.BadRequest) error {
	st := status.New(codes.InvalidArgument, "validation error")

	stWithDetails, err := st.WithDetails(details)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

func toValidatonError(vErr *validate.Error) *errdetails.BadRequest_FieldViolation {
	reason := strconv.Itoa(validationErrorReasons + int(vErr.Type))

	return &errdetails.BadRequest_FieldViolation{
		Field:       vErr.Field,
		Description: vErr.Message,
		Reason:      reason,
	}
}
