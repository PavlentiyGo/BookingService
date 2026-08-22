package main

import (
	"avitoBooking/internal/core/auth"
	core_config "avitoBooking/internal/core/config"
	core_hash "avitoBooking/internal/core/hash"
	"avitoBooking/internal/core/http/server"
	core_logger "avitoBooking/internal/core/logger"
	core_middleware "avitoBooking/internal/core/middleware"
	"avitoBooking/internal/handlers"
	"avitoBooking/internal/repository/postgres"
	"avitoBooking/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	time.Local = time.UTC
	config := core_config.NewConfigMust()

	fmt.Println(config)
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

	logger.Debug("creating hasher")
	hasher := core_hash.NewHasher(config)

	logger.Debug("creating db pool")
	pool, err := postgres.NewPool(parentCtx, config.DbConfig)
	if err != nil {
		logger.Error("failed to create db pool", zap.Error(err))
		return
	}
	txManager := postgres.NewTxManager(pool, config.DbConfig.Timeout)

	logger.Debug("creating worker")
	worker := service.NewWorker(txManager)

	ctxForWorker := core_logger.CtxWithLogger(ctx, logger)
	go worker.Start(ctxForWorker)

	conferenceService := service.NewConferenceService(true)
	serviceLayer := service.NewService(txManager, conferenceService, worker, hasher)
	handlersLayer := handlers.NewHandlers(jwtProvider, serviceLayer)

	server := core_http_server.NewServer(
		config,
		handlersLayer.GetAllRoutes(),
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
	logger.Debug("closing db pool")
	timeoutCtx, timeoutCancel := context.WithTimeout(parentCtx, config.ShutdownTimeout)
	defer timeoutCancel()
	if err = server.Shutdown(timeoutCtx); err != nil {
		logger.Error("failed to shutdown server properly", zap.Error(err))
	} else {
		logger.Debug("server was closed properly")
	}
	if err = postgres.PoolShutdown(timeoutCtx, pool); err != nil {
		logger.Error("failed to close pool properly", zap.Error(err))
	} else {
		logger.Debug("pool was closed properly")
	}
}
