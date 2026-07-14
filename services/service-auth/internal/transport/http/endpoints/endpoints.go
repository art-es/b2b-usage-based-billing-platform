package endpoints

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/register"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/resend_email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints/verify_email"
)

var (
	BindRegister                = register.Bind
	BindVerifyEmail             = verify_email.Bind
	BindResendEmailVerification = resend_email_verification.Bind
	BindLogin                   = login.Bind
)
