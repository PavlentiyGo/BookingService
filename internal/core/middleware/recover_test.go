package core_middleware

import (
	core_errors "avitoBooking/internal/core/errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RecoverSuccess(t *testing.T) {
	middleware := PanicRecoverer()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected panic")
	})
	handler := middleware(nextHandler)

	w := makeRequest(handler, "GET", "/requestId", nil, nil)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
	assert.Contains(t, w.Body.String(), core_errors.ErrInternalError.Error())
}
