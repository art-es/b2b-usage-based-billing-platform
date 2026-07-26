package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
)

var rawBodyNotImplemented, _ = json.Marshal(map[string]string{
	"message": "method is not implemented yet",
})

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

func WriteNotImplemented(w http.ResponseWriter, logger log.Logger, endpoint string) {
	logger.Log(log.Warning).
		Set("message", "called not implemeneted method").
		Set("endpoint", endpoint).
		Write()

	WriteRaw(w, http.StatusNotImplemented, rawBodyNotImplemented)
}
