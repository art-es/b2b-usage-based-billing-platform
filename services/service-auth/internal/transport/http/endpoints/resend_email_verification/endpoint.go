//go:generate mockgen -source=endpoint.go -destination=endpoint_mock_test.go -package=$GOPACKAGE
package resend_email_verification

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/resend_email_verification/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/httputil"
)

const (
	errCodeRequiredEmail = iota + 2001
	errCodeInvalidEmail
)

const (
	errMsgRequiredEmail = "Required email"
	errMsgInvalidEmail  = "Invalid email"
)

type httpRouter interface {
	Handle(pattern string, handler http.Handler)
}

type usecase interface {
	Do(ctx context.Context, email string) error
}

type handler struct {
	usecase usecase
	logger  log.Logger
}

func RegisterEndpoint(
	httpRouter httpRouter,
	usecase usecase,
	logger log.Logger,
) {
	logger = logger.Set("pkg", "internal/transport/http/endpoints/resend_email_verification")

	httpRouter.Handle("POST /v1/auth/email/resend-verification", &handler{
		usecase: usecase,
		logger:  logger,
	})
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rb requestBody

	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		httputil.WriteInvalidRequest(w)
		return
	}

	if msg, code, ok := rb.validate(); !ok {
		httputil.Write(w, http.StatusBadRequest, &httputil.BadRequestBody{
			Message: msg,
			Code:    code,
		})
		return
	}

	err = h.usecase.Do(r.Context(), rb.Email)
	if err != nil {
		if errors.Is(err, dto.ErrInvalidEmail) {
			httputil.Write(w, http.StatusBadRequest, &httputil.BadRequestBody{
				Code:    errCodeInvalidEmail,
				Message: errMsgInvalidEmail,
			})
			return
		}

		h.logger.Log(log.Error).
			Set("message", "unexpected usecase error").
			Set("error", err.Error()).
			Write()

		httputil.WriteInternalError(w)
		return
	}

	httputil.Write(w, http.StatusNoContent, nil)
}
