//go:generate mockgen -source=endpoint.go -destination=endpoint_mock_test.go -package=$GOPACKAGE
package register

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/httputil"
)

const (
	errCodeRequiredName = iota + 2001
	errCodeRequiredEmail
	errCodeRequiredPassword
	errCodeInvalidName
	errCodeInvalidEmail
	errCodeInvalidPassword
	errCodeEmailInUse
)

const (
	errMsgRequiredName     = "Required name"
	errMsgRequiredEmail    = "Required email"
	errMsgRequiredPassword = "Required password"
	errMsgInvalidName      = "Invalid name"
	errMsgInvalidEmail     = "Invalid email"
	errMsgInvalidPassword  = "Invalid password"
	errMsgEmailInUse       = "Email is already in use"
)

type httpRouter interface {
	Handle(pattern string, handler http.Handler)
}

type usecase interface {
	Do(ctx context.Context, req *dto.Request) error
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
	logger = logger.Set("pkg", "internal/transport/http/endpoints/register")

	httpRouter.Handle("POST /v1/auth/register", &handler{
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

	err = h.usecase.Do(r.Context(), &dto.Request{
		Name:     rb.Name,
		Email:    rb.Email,
		Password: rb.Password,
	})
	if err != nil {
		if errors.Is(err, dto.ErrEmailInUse) {
			httputil.Write(w, http.StatusBadRequest, &httputil.BadRequestBody{
				Code:    errCodeEmailInUse,
				Message: errMsgEmailInUse,
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

	httputil.Write(w, http.StatusAccepted, nil)
}
