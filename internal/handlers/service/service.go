package handlers_service

import (
	"avitoBooking/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	authService
	bookingService
	roomsService
}

type authService interface {
}
type bookingService interface {
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
