package endpoints

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/register"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/resend_email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/verify_email"
)

var (
	RegisterRegisterEndpoint                = register.RegisterEndpoint
	RegisterVerifyEmailEndpoint             = verify_email.RegisterEndpoint
	RegisterResendEmailVerificationEndpoint = resend_email_verification.RegisterEndpoint
	RegisterLoginEndpoint                   = login.RegisterEndpoint
)
