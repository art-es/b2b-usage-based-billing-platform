package usecases

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/finish_all_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_me"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/refresh_session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/resend_email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/verify_email"
)

var (
	NewRegisterUsecase                = register.NewUsecase
	NewVerifyEmailUsecase             = verify_email.NewUsecase
	NewResendEmailVerificationUsecase = resend_email_verification.NewUsecase
	NewLoginUsecase                   = login.NewUsecase
	NewRefreshSessionUsecase          = refresh_session.NewUsecase
	NewGetSessionsUsecase             = get_sessions.NewUsecase
	NewFinishAllSessionsUsecase       = finish_all_sessions.NewUsecase
	NewGetMeUsecase                   = get_me.NewUsecase
)
