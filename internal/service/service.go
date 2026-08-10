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
	worker            WorkerInterface // TODO INTERFACE
}

func NewService(
	txManager repository.TxManager,
	conferenceService *ConferenceService,
	worker WorkerInterface,
) *Service {

	service := &Service{
		txManager:         txManager,
		conferenceService: conferenceService,
		worker:            worker,
	}

	storage, ok := txManager.(repository.Storage)

	if ok {
		service.roomsRepo = storage.GetRoomsRepo()
		service.bookingRepo = storage.GetBookingRepo()
	}
	return service
}
