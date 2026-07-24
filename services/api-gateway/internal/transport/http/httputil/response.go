package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
)

var (
	bodyUnauthorized   = []byte(`{"message":"Unauthorized"}`)
	bodyInvalidRequest = []byte(`{"code":1001,"message":"Invalid request format"}`)
	bodyInternalError  = []byte(`{"message":"Internal error"}`)
	bodyNotImplemented = []byte(`{"message":"Method is not implemented yet"}`)
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

func WriteInternalError(w http.ResponseWriter, logger log.Logger, err error) {
	logObj := logger.Log(log.Error).
		Set("message", "unexpected error")

	if err != nil {
		logObj.Set("error", err.Error())
	}

	logObj.Write()

	WriteRaw(w, http.StatusInternalServerError, bodyInternalError)
}

func WriteNotImplemented(w http.ResponseWriter, logger log.Logger, method string) {
	logger.Log(log.Warning).
		Set("message", "called not implemeneted method").
		Set("method", method).
		Write()

	WriteRaw(w, http.StatusNotImplemented, bodyNotImplemented)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	}

	Write(w, http.StatusBadRequest, &BadRequestBody{Message: msg})
}
