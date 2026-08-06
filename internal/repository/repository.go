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
}
type BookingRepository interface {
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
}
type ScheduleRepository interface {
}
type SlotsRepository interface {
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
}
