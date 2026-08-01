package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	uerrors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/ptr"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/httputil"
)

var rawBodyInvalidRequestFormat, _ = json.Marshal(&openapi.BadRequestResponse{
	Message: ptr.To("invalid request format"),
})

var rawBodyInternalError, _ = json.Marshal(openapi.InternalErrorResponse{
	Message: ptr.To("internal error"),
})

func getRequestErrorHandlerFunc() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, _ *http.Request, _ error) {
		httputil.WriteRaw(w, http.StatusBadRequest, rawBodyInvalidRequestFormat)
	}
}

func getResponseErrorHandlerFunc(logger log.Logger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			logger.Log(log.Error).
				Set("message", "unexpected call of openapi ResponseErrorHandlerFunc").
				Set("endpoint", endpointByRequest(r)).
				Write()

			httputil.Write(w, http.StatusInternalServerError, rawBodyInternalError)
			return
		}

		if errors.Is(err, errUnimplemented) {
			httputil.WriteNotImplemented(w, logger, endpointByRequest(r))
			return
		}

		if errors.Is(err, app_errors.ErrUnauthorized) {
			httputil.WriteUnauthorized(w)
			return
		}

		{
			var vErr *app_errors.ValidationError
			if errors.As(err, &vErr) {
				httputil.Write(w, http.StatusBadRequest, &openapi.BadRequestResponse{
					Errors: &[]openapi.BadRequestResponseError{
						convertValidationError(vErr),
					},
				})
				return
			}
		}

		{
			var vErrs app_errors.ValidationErrors
			if errors.As(err, &vErrs) {
				errs := []openapi.BadRequestResponseError{}
				for _, uErr := range vErrs {
					errs = append(errs, convertValidationError(uErr))
				}

				httputil.Write(w, http.StatusBadRequest, &openapi.BadRequestResponse{
					Errors: &errs,
				})
				return
			}
		}

		logger.Log(log.Error).
			Set("message", "unexpected endpoint error").
			Set("endpoint", endpointByRequest(r)).
			Set("error", err.Error()).
			Write()

		httputil.Write(w, http.StatusInternalServerError, rawBodyInternalError)
	}
}

func convertValidationError(uErr *uerrors.ValidationError) openapi.BadRequestResponseError {
	var code openapi.BadRequestResponseErrorCode
	switch uErr.Code {
	case 1001:
		code = openapi.N1001
	case 1002:
		code = openapi.N1002
	case 1003:
		code = openapi.N1003
	case 1004:
		code = openapi.N1004
	case 1005:
		code = openapi.N1005
	case 1006:
		code = openapi.N1006
	}

	return openapi.BadRequestResponseError{
		Code:    ptr.ToOrNil(code),
		Field:   ptr.ToOrNil(uErr.Field),
		Message: ptr.ToOrNil(uErr.Message),
	}
}

func endpointByRequest(r *http.Request) string {
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
}
