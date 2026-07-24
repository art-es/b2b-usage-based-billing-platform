package post_v1_auth_login

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/httputil"
)

type authService interface {
	Login(ctx context.Context, req *auth_service.LoginRequest) (*auth_service.LoginResponse, error)
}

type requestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type responseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Handler struct {
	logger      log.Logger
	authService authService
}

func NewHandler(
	logger log.Logger,
	authService authService,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/http/openapi/handlers/post_v1_auth_login")

	return &Handler{
		logger:      logger,
		authService: authService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var rb requestBody

	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		httputil.WriteInvalidRequest(w)
		return
	}

	res, err := h.authService.Login(r.Context(), &auth_service.LoginRequest{
		Email:    rb.Email,
		Password: rb.Password,
	})
	if err != nil {
		httputil.WriteInternalError(w, h.logger, err)
		return
	}

	if res.RequestError != nil {
		httputil.Write(w, http.StatusBadRequest, &httputil.BadRequestBody{
			Code:    res.RequestError.Code,
			Message: res.RequestError.Message,
		})
		return
	}

	httputil.Write(w, http.StatusOK, &responseBody{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
