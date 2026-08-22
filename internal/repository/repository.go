package repository

import (
	"avitoBooking/internal/core/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository interface {
	Register(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (domain.User, error)
}
type BookingRepository interface {
	CreateBooking(
		ctx context.Context,
		booking domain.BookingRequest,
	) (domain.BookingRequest, error)
	GetSlotById(
		ctx context.Context,
		slotId uuid.UUID,
	) (domain.Slot, error)
	ChangeBookingStatus(
		ctx context.Context,
		bookingId uuid.UUID,
		newStatus domain.BookingStatus,
	) (domain.BookingRequest, error)
	GetBookingsByUserId(
		ctx context.Context,
		userId uuid.UUID,
	) ([]domain.BookingRequest, error)
	GetBookingsCount(
		ctx context.Context,
	) (int, error)
	GetBookingsWithPagination(
		ctx context.Context,
		pagination domain.Pagination,
	) ([]domain.BookingRequest, error)
}
type RoomsRepository interface {
	CreateRoom(
		ctx context.Context,
		room domain.Room,
	) (domain.Room, error)
	GetRooms(
		ctx context.Context,
	) ([]domain.Room, error)
	CreateSchedule(
		ctx context.Context,
		roomSchedule domain.RoomSchedule,
	) (domain.RoomSchedule, error)
	GetSlots(
		ctx context.Context,
		roomId uuid.UUID,
		dateTime time.Time,
	) ([]domain.Slot, error)

	GetRoomsSchedule(
		ctx context.Context,
	) ([]domain.RoomSchedule, error)
}
type WorkerRepository interface {
	CreateSlots(
		ctx context.Context,
		slots []domain.Slot,
	) error
	CleanSlots(
		ctx context.Context,
	) error
}

type Executor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
type TxManager interface {
	WithinTx(
		ctx context.Context,
		fn func(txCtx context.Context) error,
	) error
	GetExecutor(ctx context.Context) Executor
	GetTimeout() time.Duration
}
type Storage interface {
	GetRoomsRepo() RoomsRepository
	GetBookingRepo() BookingRepository
	GetWorkerRepo() WorkerRepository
	GetAuthRepo() AuthRepository
}
