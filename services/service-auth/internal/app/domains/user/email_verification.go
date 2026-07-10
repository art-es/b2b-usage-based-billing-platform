package user

import "fmt"

type EmailVerification struct {
	Token  string
	Email  string
	UserID string
}

func (v *EmailVerification) EmailSubject() string {
	return "Verify your new account"
}

func (v *EmailVerification) EmailContent() string {
	return fmt.Sprintf(
		"To verify your email address, please follow the link: https://example.com/verify-email/%s",
		v.Token,
	)
}
