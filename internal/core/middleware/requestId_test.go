package core_middleware

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_RequestIdEmpty(t *testing.T) {

	middleware := RequestId()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware(nextHandler)

	w := makeRequest(handler, "GET", "/requestId", nil, nil)

	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	assert.Equal(t, w.Code, http.StatusOK)
}
func Test_RequestIdExist(t *testing.T) {
	middleware := RequestId()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware(nextHandler)

	headers := make(map[string]string, 1)
	uuidRequest := uuid.NewString()
	headers["X-Request-ID"] = uuidRequest
	w := makeRequest(handler, "GET", "/requestId", nil, headers)

	assert.Equal(t, w.Header().Get("X-Request-ID"), uuidRequest)
	assert.Equal(t, w.Code, http.StatusOK)
}
