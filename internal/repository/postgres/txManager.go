package postgres

import (
	core_logger "avitoBooking/internal/core/logger"
	"avitoBooking/internal/repository"
	postgres_repository "avitoBooking/internal/repository/postgres/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TxManager struct {
	pool           *pgxpool.Pool
	requestTimeout time.Duration
	authRepo       repository.AuthRepository
	bookingRepo    repository.BookingRepository
	roomsRepo      repository.RoomsRepository
	workerRepo     repository.WorkerRepository
}
type txCtxKey struct{}

func NewTxManager(
	pool *pgxpool.Pool,
	timeout time.Duration,
) *TxManager {
	txManager := &TxManager{
		pool:           pool,
		requestTimeout: timeout,
	}
	txManager.authRepo = postgres_repository.NewAuthRepository(txManager)
	txManager.roomsRepo = postgres_repository.NewRoomsRepository(txManager)
	txManager.bookingRepo = postgres_repository.NewBookingRepository(txManager)
	txManager.workerRepo = postgres_repository.NewWorkerRepository(txManager)
	return txManager
}

func (m *TxManager) WithinTx(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) error {
	tx, err := m.pool.Begin(ctx)
	defer func() {
		if err = tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger := core_logger.FromContext(ctx)
			logger.Error("failed to rollback transaction", zap.Error(err))
		}
	}()
	txCtx := context.WithValue(ctx, txCtxKey{}, tx)
	if err = fn(txCtx); err != nil {
		return fmt.Errorf("error during transaction: %w", err)
	}
	if err = tx.Commit(txCtx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	return nil
}
func (m *TxManager) GetExecutor(ctx context.Context) repository.Executor {
	if tx, ok := ctx.Value(txCtxKey{}).(pgx.Tx); ok {
		return tx
	}
	return m.pool
}
func (m *TxManager) GetTimeout() time.Duration {
	return m.requestTimeout
}
func (m *TxManager) GetRoomsRepo() repository.RoomsRepository {
	return m.roomsRepo
}
func (m *TxManager) GetBookingRepo() repository.BookingRepository {
	return m.bookingRepo
}
func (m *TxManager) GetWorkerRepo() repository.WorkerRepository { return m.workerRepo }
func (m *TxManager) GetAuthRepo() repository.AuthRepository {
	return m.authRepo
}
