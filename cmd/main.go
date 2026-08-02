package main

import (
	"avitoBooking/internal/core/auth"
	core_config "avitoBooking/internal/core/config"
	"avitoBooking/internal/core/http/server"
	core_logger "avitoBooking/internal/core/logger"
	core_middleware "avitoBooking/internal/core/middleware"
	core_routes "avitoBooking/internal/core/routes"
	"avitoBooking/internal/handlers"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	config := core_config.NewConfigMust()

	logger, err := core_logger.NewLogger(config)
	if err != nil {
		log.Fatal("failed to create new logger", err)
	}
	defer logger.Close()

	parentCtx := context.Background()

	ctx, cancel := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger.Debug("creating jwt provider")
	jwtProvider := auth.NewJwtConfig(config.LifeTime, config.SigningKey)

	handlers := handlers.NewHandlers(jwtProvider)
	var routes []core_routes.Route
	routes = append(routes, handlers.AuthRoutes()...)
	routes = append(routes, handlers.SlotsRoutes()...)
	routes = append(routes, handlers.ScheduleRoutes()...)
	routes = append(routes, handlers.BookingRoutes()...)
	routes = append(routes, handlers.RoomsRoutes()...)

	server := core_http_server.NewServer(
		config,
		routes,
		core_middleware.RequestId(),
		core_middleware.Logger(logger),
		core_middleware.Trace(),
		core_middleware.PanicRecoverer(),
	)

	logger.Debug("starting server", zap.String("Address", server.Addr))
	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", zap.Error(err))
		}
	}()
	<-ctx.Done()
	logger.Debug("closing server")
	timeoutCtx, timeoutCancel := context.WithTimeout(parentCtx, time.Second*5)
	defer timeoutCancel()

	if err = server.Shutdown(timeoutCtx); err != nil {
		logger.Error("failed to shutdown server properly", zap.Error(err))
	}
}
