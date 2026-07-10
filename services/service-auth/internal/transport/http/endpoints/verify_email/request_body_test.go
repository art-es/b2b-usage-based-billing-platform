package verify_email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBody_Validate(t *testing.T) {
	t.Run("empty token", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Token: "",
		}).validate()

		assert.Equal(t, "Required token", msg)
		assert.Equal(t, 2001, code)
		assert.False(t, ok)
	})

	t.Run("invalid token", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Token: "invalid-token",
		}).validate()

		assert.Equal(t, "Invalid token", msg)
		assert.Equal(t, 2002, code)
		assert.False(t, ok)
	})

	t.Run("ok", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Token: "00000000-0000-0000-0000-000000000001",
		}).validate()

		assert.Empty(t, msg)
		assert.Empty(t, code)
		assert.True(t, ok)
	})
}
