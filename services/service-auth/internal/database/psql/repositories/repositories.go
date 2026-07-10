package repositories

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories/verification"
)

var (
	NewUserRepository         = user.NewRepository
	NewVerificationRepository = verification.NewRepository
)
