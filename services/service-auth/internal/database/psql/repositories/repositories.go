package repositories

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories/email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories/sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories/user"
)

var (
	NewUserRepository              = user.NewRepository
	NewEmailVerificationRepository = email_verification.NewRepository
	NewSessionsRepository          = sessions.NewRepository
)
