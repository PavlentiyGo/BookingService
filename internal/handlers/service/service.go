package handlers_service

import (
	"avitoBooking/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --all --output=./mocks --outpkg=mocks
type BookingService interface {
	authService
	bookService
	roomsService
}

type authService interface {
	Register(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
	Login(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}
type bookService interface {
	CreateBooking(
		ctx context.Context,
		booking domain.BookingRequest,
	) (domain.BookingRequest, error)
	GetMyBookings(
		ctx context.Context,
		userId uuid.UUID,
	) ([]domain.BookingRequest, error)
	CancelUserBooking(
		ctx context.Context,
		userId uuid.UUID,
		bookingId string,
	) (domain.BookingRequest, error)
	GetBookings(
		ctx context.Context,
		pagination domain.Pagination,
	) ([]domain.BookingRequest, int, error)
}
type roomsService interface {
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
		date string,
	) ([]domain.Slot, error)
}
