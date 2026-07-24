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
	br := &errdetails.BadRequest{}
	for _, vErr := range vErrs {
		br.FieldViolations = append(br.FieldViolations, toValidatonError(vErr))
	}

	st := status.New(codes.InvalidArgument, "validation error")
	st.WithDetails(br)

	return st.Err()
}

func ValidationError(vErr *validate.Error) error {
	if vErr == nil {
		return nil
	}

	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, toValidatonError(vErr))

	st := status.New(codes.InvalidArgument, "validation error")
	st.WithDetails(br)

	return st.Err()
}

func toValidatonError(vErr *validate.Error) *errdetails.BadRequest_FieldViolation {
	reason := strconv.Itoa(validationErrorReasons + int(vErr.Type))

	return &errdetails.BadRequest_FieldViolation{
		Field:       vErr.Field,
		Description: vErr.Message,
		Reason:      reason,
	}
}
