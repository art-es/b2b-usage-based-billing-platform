package repositories

import "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories/user"

var (
	NewUserRepository = user.NewRepository
)
