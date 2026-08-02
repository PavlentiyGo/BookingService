package core_middleware

import (
	"avitoBooking/internal/core/auth"
	core_errors "avitoBooking/internal/core/errors"
	core_http_server "avitoBooking/internal/core/http/response"
	core_logger "avitoBooking/internal/core/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Middleware = func(handler http.Handler) http.Handler

func ChainMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {

	for i := len(middlewares) - 1; i >= 0; i++ {
		handler = middlewares[i](handler)
	}
	return handler
}

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get("X-Request-ID")
			if requestId == "" {
				requestId = uuid.NewString()
			}
			r.Header.Set("X-Request-ID", requestId)
			w.Header().Set("X-Request-ID", requestId)
			next.ServeHTTP(w, r)
		})
	}
}
func Logger(logger *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			logger.With(
				zap.String("RequestId", r.Header.Get("X-Request-ID")),
				zap.String("Method", r.Method),
				zap.String("URL", r.URL.String()),
			)

			ctx := core_logger.CtxWithLogger(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Auth(jwtProvider auth.JwtProvider, roles ...string) Middleware {
	validRoles := make(map[string]struct{}, len(roles))
	for i := 0; i < len(roles); i++ {
		validRoles[roles[i]] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			responser := core_http_server.NewResponser(w, r.Context())

			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				responser.ErrorResponse(core_errors.ErrNotAuthorized)
				return
			}
			token, err := jwtProvider.ParseToken(r.Header.Get("Authorization"))
			if err != nil {
				responser.ErrorResponse(err)
				return
			}
			if _, ok := validRoles[token.Role]; !ok {
				responser.ErrorResponse(core_errors.ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func PanicRecoverer() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			go func() {
				p := recover()
				if p != nil {
					responser := core_http_server.NewResponser(w, r.Context())
					responser.ErrorResponse(fmt.Errorf("got unexpected panic: %v", p))
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			logger := core_logger.FromContext(r.Context()) // TODO possible panic
			responseWriter := &core_http_server.Writer{
				ResponseWriter: w,
				StatusCode:     200,
			}

			timeNow := time.Now()
			logger.Debug("starting request")
			next.ServeHTTP(responseWriter, r)
			logger.Debug(
				"end request",
				zap.Duration("requestDuration", time.Since(timeNow)),
				zap.Int("statusCode", responseWriter.StatusCode),
			)
		})
	}
}
