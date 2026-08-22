package handlers

import (
	mocks2 "avitoBooking/internal/core/auth/mocks"
	core_http_server "avitoBooking/internal/core/http/server"
	core_logger "avitoBooking/internal/core/logger"
	"avitoBooking/internal/handlers/service/mocks"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

type testDeps struct {
	service     *mocks.BookingService
	router      *http.ServeMux
	jwtProvider *mocks2.JwtProvider
}

func newTestDeps(t *testing.T) *testDeps {
	service := mocks.NewBookingService(t)
	jwtProvider := mocks2.NewJwtProvider(t)

	handler := NewHandlers(jwtProvider, service)
	mux := core_http_server.ChainRoutes(handler.GetAllRoutes()...)

	return &testDeps{
		service:     service,
		router:      mux,
		jwtProvider: jwtProvider,
	}
}
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
