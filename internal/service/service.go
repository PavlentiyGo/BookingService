package service

import (
	"avitoBooking/internal/repository"
)

type Service struct {
	txManager    repository.TxManager
	authRepo     repository.AuthRepository
	bookingRepo  repository.BookingRepository
	roomsRepo    repository.RoomsRepository
	scheduleRepo repository.ScheduleRepository
	slotsRepo    repository.SlotsRepository
}

func NewService(
	txManager repository.TxManager,
) *Service {

	service := &Service{
		txManager: txManager,
	}

	storage, ok := txManager.(repository.Storage)

	if ok {
		service.roomsRepo = storage.GetRoomsRepo()
		service.bookingRepo = storage.GetBookingRepo()
	}
	return service
}
