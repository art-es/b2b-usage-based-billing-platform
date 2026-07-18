//go:generate mockgen -source=endpoint.go -destination=endpoint_mock_test.go -package=$GOPACKAGE
package verify_email

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/verify_email/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/httputil"
)

const (
	errCodeRequiredToken = iota + 2001
	errCodeInvalidToken
)

const (
	errMsgRequiredToken = "Required token"
	errMsgInvalidToken  = "Invalid token"
)

type router interface {
	Handle(pattern string, handler http.Handler)
}

type usecase interface {
	Do(ctx context.Context, token string) error
}

type handler struct {
	usecase usecase
	logger  log.Logger
}

func Bind(
	router router,
	usecase usecase,
	logger log.Logger,
) {
	logger = logger.Set("pkg", "internal/transport/http/endpoints/verify_email")

	router.Handle("POST /v1/auth/email/verify", &handler{
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

	err = h.usecase.Do(r.Context(), rb.Token)
	if err != nil {
		if errors.Is(err, dto.ErrInvalidToken) {
			httputil.Write(w, http.StatusBadRequest, &httputil.BadRequestBody{
				Code:    errCodeInvalidToken,
				Message: errMsgInvalidToken,
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
