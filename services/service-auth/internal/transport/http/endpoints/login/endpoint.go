//go:generate mockgen -source=endpoint.go -destination=endpoint_mock_test.go -package=$GOPACKAGE
package login

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/httputil"
)

const (
	errCodeRequiredEmail = iota + 2001
	errCodeRequiredPassword
	errCodeInvalidEmail
	errCodeInvalidPassword
	errCodeWrongCredentials
	errCodeEmailNotVerified
)

const (
	errMsgRequiredEmail    = "Required email"
	errMsgRequiredPassword = "Required password"
	errMsgInvalidEmail     = "Invalid email"
	errMsgInvalidPassword  = "Invalid password"
	errMsgWrongCredentials = "Wrong credentials"
	errMsgEmailNotVerified = "Email not verified"
)

type httpRouter interface {
	Handle(pattern string, handler http.Handler)
}

type usecase interface {
	Do(ctx context.Context, req *dto.Request) (*dto.Response, error)
}

type responseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	logger = logger.Set("pkg", "internal/transport/http/endpoints/login")

	httpRouter.Handle("POST /v1/auth/login", &handler{
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

	res, err := h.usecase.Do(r.Context(), &dto.Request{
		Email:    rb.Email,
		Password: rb.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrWrongCredentials):
			httputil.Write(w, http.StatusBadRequest, &httputil.BadRequestBody{
				Code:    errCodeWrongCredentials,
				Message: errMsgWrongCredentials,
			})
		case errors.Is(err, dto.ErrEmailNotVerified):
			httputil.Write(w, http.StatusBadRequest, &httputil.BadRequestBody{
				Code:    errCodeEmailNotVerified,
				Message: errMsgEmailNotVerified,
			})
		default:
			h.logger.Log(log.Error).
				Set("message", "unexpected usecase error").
				Set("error", err.Error()).
				Write()

			httputil.WriteInternalError(w)
		}
		return
	}

	httputil.Write(w, http.StatusOK, &responseBody{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
