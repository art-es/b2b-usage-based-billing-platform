package httputil

import (
	"encoding/json"
	"net/http"
)

const (
	ErrCodeInvalidRequest = iota + 1001
)

var (
	bodyInvalidRequest = []byte(`{"code":1001,"message":"Invalid request format"}`)
	bodyInternalError  = []byte(`{"message":"Internal error"}`)
)

type BadRequestBody struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func WriteRaw(w http.ResponseWriter, code int, body []byte) {
	w.WriteHeader(code)

	if body != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func Write(w http.ResponseWriter, code int, body any) {
	var encBody []byte
	if body != nil {
		encBody, _ = json.Marshal(body)
	}

	WriteRaw(w, code, encBody)
}

func WriteInvalidRequest(w http.ResponseWriter) {
	WriteRaw(w, http.StatusBadRequest, bodyInvalidRequest)
}

func WriteInternalError(w http.ResponseWriter) {
	WriteRaw(w, http.StatusInternalServerError, bodyInternalError)
}
