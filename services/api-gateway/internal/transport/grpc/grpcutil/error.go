package grpcutil

import (
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
)

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return handleInvalidArgument(st)
	default:
		return err
	}
}

func handleInvalidArgument(st *status.Status) error {
	var vErrs app_errors.ValidationErrors

	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range d.FieldViolations {
				code, _ := strconv.Atoi(violation.Reason)

				vErrs = append(vErrs, &app_errors.ValidationError{
					Field:   violation.Field,
					Message: violation.Description,
					Code:    code,
				})
			}
		}
	}

	if len(vErrs) == 0 {
		return st.Err()
	}

	return vErrs
}
