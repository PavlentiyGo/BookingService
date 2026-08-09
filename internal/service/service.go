package service

import (
	"avitoBooking/internal/repository"
)

type Service struct {
	txManager         repository.TxManager
	authRepo          repository.AuthRepository
	bookingRepo       repository.BookingRepository
	roomsRepo         repository.RoomsRepository
	conferenceService *ConferenceService
}

func NewService(
	txManager repository.TxManager,
	conferenceService *ConferenceService,
) *Service {

	service := &Service{
		txManager:         txManager,
		conferenceService: conferenceService,
	}

	storage, ok := txManager.(repository.Storage)

	if ok {
		service.roomsRepo = storage.GetRoomsRepo()
		service.bookingRepo = storage.GetBookingRepo()
	}
	return service
}
