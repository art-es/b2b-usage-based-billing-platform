package hmac_sha256

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	svc := NewService()

	t.Run("same string, same secret", func(t *testing.T) {
		fooHash1, err := svc.Generate([]byte("a"), "foo")
		require.NoError(t, err)

		fooHash2, err := svc.Generate([]byte("a"), "foo")
		require.NoError(t, err)

		assert.Equal(t, fooHash1, fooHash2)
	})

	t.Run("same string, different secrets", func(t *testing.T) {
		fooHash1, err := svc.Generate([]byte("a"), "foo")
		require.NoError(t, err)

		fooHash2, err := svc.Generate([]byte("b"), "foo")
		require.NoError(t, err)

		assert.NotEqual(t, fooHash1, fooHash2)
	})

	t.Run("different string", func(t *testing.T) {
		fooHash, err := svc.Generate([]byte("a"), "foo")
		require.NoError(t, err)

		barHash, err := svc.Generate([]byte("a"), "bar")
		require.NoError(t, err)

		assert.NotEqual(t, fooHash, barHash)
	})
}
