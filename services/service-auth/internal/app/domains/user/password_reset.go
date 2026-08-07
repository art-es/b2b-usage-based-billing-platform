package user

import "fmt"

type PasswordReset struct {
	Token  string
	UserID string
	Email  string
}

func (v *PasswordReset) EmailSubject() string {
	return "Reset password"
}

func (v *PasswordReset) EmailContent() string {
	return fmt.Sprintf(
		"To reset password, please follow the link: https://example.com/reset-password/%s",
		v.Token,
	)
}
