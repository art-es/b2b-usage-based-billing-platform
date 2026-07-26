package grpcutil

import (
	"strconv"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/uerrors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	var vErrs uerrors.ValidationErrors

	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range d.FieldViolations {
				code, _ := strconv.Atoi(violation.Reason)

				vErrs = append(vErrs, &uerrors.ValidationError{
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
