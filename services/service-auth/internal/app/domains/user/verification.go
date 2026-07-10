package user

import "fmt"

type Verification struct {
	Token string
	Email string
}

func (v *Verification) EmailSubject() string {
	return "Verify your new account"
}

func (v *Verification) EmailContent() string {
	return fmt.Sprintf(
		"To verify your email address, please follow the link: https://example.com/verify-email/%s",
		v.Token,
	)
}
