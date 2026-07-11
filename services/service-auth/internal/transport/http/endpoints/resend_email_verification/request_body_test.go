package resend_email_verification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBody_Validate(t *testing.T) {
	t.Run("empty email", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Email: "",
		}).validate()

		assert.Equal(t, "Required email", msg)
		assert.Equal(t, 2001, code)
		assert.False(t, ok)
	})

	t.Run("email is too long", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Email: string(make([]byte, 101)),
		}).validate()

		assert.Equal(t, "Invalid email", msg)
		assert.Equal(t, 2002, code)
		assert.False(t, ok)
	})

	t.Run("invalid email", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Email: "invalid-email",
		}).validate()

		assert.Equal(t, "Invalid email", msg)
		assert.Equal(t, 2002, code)
		assert.False(t, ok)
	})

	t.Run("ok", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Email: "a@a.a",
		}).validate()

		assert.Empty(t, msg)
		assert.Empty(t, code)
		assert.True(t, ok)
	})
}
