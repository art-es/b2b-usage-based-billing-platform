package register

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBody_Validate(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "",
			Email:    "",
			Password: "",
		}).validate()

		assert.Equal(t, "Required name", msg)
		assert.Equal(t, 2001, code)
		assert.False(t, ok)
	})

	t.Run("empty email", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "a",
			Email:    "",
			Password: "",
		}).validate()

		assert.Equal(t, "Required email", msg)
		assert.Equal(t, 2002, code)
		assert.False(t, ok)
	})

	t.Run("empty password", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "a",
			Email:    "a",
			Password: "",
		}).validate()

		assert.Equal(t, "Required password", msg)
		assert.Equal(t, 2003, code)
		assert.False(t, ok)
	})

	t.Run("name length is too short", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "a",
			Email:    "a",
			Password: "a",
		}).validate()

		assert.Equal(t, "Invalid name", msg)
		assert.Equal(t, 2004, code)
		assert.False(t, ok)
	})

	t.Run("name length is too long", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "a" + string(make([]byte, 100)),
			Email:    "a",
			Password: "a",
		}).validate()

		assert.Equal(t, "Invalid name", msg)
		assert.Equal(t, 2004, code)
		assert.False(t, ok)
	})

	t.Run("name's first letter is not letter", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "1aa",
			Email:    "a",
			Password: "a",
		}).validate()

		assert.Equal(t, "Invalid name", msg)
		assert.Equal(t, 2004, code)
		assert.False(t, ok)
	})

	t.Run("email is too long", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "aa1",
			Email:    string(make([]byte, 101)),
			Password: "a",
		}).validate()

		assert.Equal(t, "Invalid email", msg)
		assert.Equal(t, 2005, code)
		assert.False(t, ok)
	})

	t.Run("email is invalid", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "aa1",
			Email:    "aa1",
			Password: "a",
		}).validate()

		assert.Equal(t, "Invalid email", msg)
		assert.Equal(t, 2005, code)
		assert.False(t, ok)
	})

	t.Run("password is too short", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "aa1",
			Email:    "aa1@a.a",
			Password: "a",
		}).validate()

		assert.Equal(t, "Invalid password", msg)
		assert.Equal(t, 2006, code)
		assert.False(t, ok)
	})

	t.Run("password is too long", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "aa1",
			Email:    "aa1@a.a",
			Password: string(make([]byte, 65)),
		}).validate()

		assert.Equal(t, "Invalid password", msg)
		assert.Equal(t, 2006, code)
		assert.False(t, ok)
	})

	t.Run("password with no digits", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "aa1",
			Email:    "aa1@a.a",
			Password: "aaaaaaaa",
		}).validate()

		assert.Equal(t, "Invalid password", msg)
		assert.Equal(t, 2006, code)
		assert.False(t, ok)
	})

	t.Run("password with no letters", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "aa1",
			Email:    "aa1@a.a",
			Password: "11111111",
		}).validate()

		assert.Equal(t, "Invalid password", msg)
		assert.Equal(t, 2006, code)
		assert.False(t, ok)
	})

	t.Run("ok", func(t *testing.T) {
		msg, code, ok := (&requestBody{
			Name:     "aa1",
			Email:    "aa1@a.a",
			Password: "aaaa1111",
		}).validate()

		assert.Empty(t, msg)
		assert.Empty(t, code)
		assert.True(t, ok)
	})
}
