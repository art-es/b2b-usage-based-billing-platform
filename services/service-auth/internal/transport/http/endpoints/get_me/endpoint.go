package get_me

import (
	"context"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_me/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/httputil"
)

type router interface {
	Handle(pattern string, handler http.Handler)
}

type authorizer interface {
	Authorize(w http.ResponseWriter, r *http.Request) (*jwt.Claims, bool)
}

type usecase interface {
	Do(ctx context.Context, claims *jwt.Claims) (*dto.Response, error)
}

type responseBody struct {
	SessionID string            `json:"session_id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Orgn      *responseBodyOrgn `json:"orgn,omitempty"`
}

type responseBodyOrgn struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type handler struct {
	authorizer authorizer
	usecase    usecase
	logger     log.Logger
}

func Bind(
	router router,
	authorizer authorizer,
	usecase usecase,
	logger log.Logger,
) {
	logger = logger.Set("pkg", "internal/transport/http/endpoints/get_me")

	router.Handle("POST /v1/auth/login", &handler{
		authorizer: authorizer,
		usecase:    usecase,
		logger:     logger,
	})
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtc, ok := h.authorizer.Authorize(w, r)
	if !ok {
		return
	}

	res, err := h.usecase.Do(r.Context(), jwtc)
	if err != nil {
		h.logger.Log(log.Error).
			Set("message", "unexpected usecase error").
			Set("error", err.Error()).
			Write()

		httputil.WriteInternalError(w)
		return
	}

	httputil.Write(w, http.StatusOK, &responseBody{
		SessionID: res.SessionID,
		Name:      res.User.Name,
		Email:     res.User.Email,
		Orgn:      convertOrgn(res.Orgn),
	})
}

func convertOrgn(in *dto.ResponseOrgn) *responseBodyOrgn {
	if in != nil {
		return ptr.To(responseBodyOrgn(*in))
	}
	return nil
}
