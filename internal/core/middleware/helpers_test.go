package core_middleware

import (
	core_logger "avitoBooking/internal/core/logger"
	"bytes"
	"net/http"
	"net/http/httptest"

	"go.uber.org/zap"
)

func makeRequest(
	handler http.Handler,
	method string,
	target string,
	body []byte,
	headers map[string]string,
) *httptest.ResponseRecorder {

	req := httptest.NewRequest(method, target, bytes.NewBuffer(body))

	ctx := core_logger.CtxWithLogger(req.Context(), &core_logger.Logger{
		Logger: zap.NewNop(),
	})
	for header, val := range headers {
		req.Header.Set(header, val)
	}

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req.WithContext(ctx))
	return w
}
