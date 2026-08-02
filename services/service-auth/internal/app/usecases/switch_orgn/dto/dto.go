package dto

import "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"

type Request struct {
	Auth   *jwt.Claims
	OrgnID string
}
