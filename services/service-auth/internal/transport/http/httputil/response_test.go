package httputil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteRaw(t *testing.T) {
	t.Run("body is nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		code := http.StatusInternalServerError

		assert.NotEqual(t, code, w.Code)
		assert.Empty(t, w.Header().Get("Content-Type"))
		assert.Empty(t, w.Body.String())

		WriteRaw(w, code, nil)

		assert.Equal(t, code, w.Code)
		assert.Empty(t, w.Header().Get("Content-Type"))
		assert.Empty(t, w.Body.String())
	})

	t.Run("body is not nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		code := http.StatusInternalServerError
		body := `{"message":"something happened"}`

		assert.NotEqual(t, code, w.Code)
		assert.Empty(t, w.Header().Get("Content-Type"))
		assert.NotEqual(t, body, w.Body.String())

		WriteRaw(w, code, []byte(body))

		assert.Equal(t, code, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, body, w.Body.String())
	})
}

func TestWrite(t *testing.T) {
	t.Run("body is nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		code := http.StatusInternalServerError

		assert.NotEqual(t, code, w.Code)
		assert.Empty(t, w.Header().Get("Content-Type"))
		assert.Empty(t, w.Body.String())

		Write(w, code, nil)

		assert.Equal(t, code, w.Code)
		assert.Empty(t, w.Header().Get("Content-Type"))
		assert.Empty(t, w.Body.String())
	})

	t.Run("body is not nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		code := http.StatusInternalServerError
		body := map[string]any{"message": "something happened"}

		assert.NotEqual(t, code, w.Code)
		assert.Empty(t, w.Header().Get("Content-Type"))
		assert.Empty(t, w.Body.String())

		Write(w, code, body)

		fmt.Println(w.Code, w.Body.String())
		assert.Equal(t, code, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, `{"message":"something happened"}`, w.Body.String())
	})
}

func TestWriteInvalidRequest(t *testing.T) {
	w := httptest.NewRecorder()

	WriteInvalidRequest(w)

	expCode := 400
	expBody := `{"code":1001,"message":"Invalid request format"}`

	assert.Equal(t, expCode, w.Code)
	assert.JSONEq(t, expBody, w.Body.String())
}

func TestWriteInternalError(t *testing.T) {
	w := httptest.NewRecorder()

	WriteInternalError(w)

	expCode := 500
	expBody := `{"message":"Internal error"}`

	assert.Equal(t, expCode, w.Code)
	assert.JSONEq(t, expBody, w.Body.String())
}
