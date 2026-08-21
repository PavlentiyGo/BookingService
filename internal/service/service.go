package service

import (
	core_hash "avitoBooking/internal/core/hash"
	"avitoBooking/internal/repository"
)

type Service struct {
	txManager         repository.TxManager
	authRepo          repository.AuthRepository
	bookingRepo       repository.BookingRepository
	roomsRepo         repository.RoomsRepository
	conferenceService *ConferenceService
	worker            WorkerInterface
	hasher            core_hash.HasherInterface
}

func NewService(
	txManager repository.TxManager,
	conferenceService *ConferenceService,
	worker WorkerInterface,
	hasher core_hash.HasherInterface,
) *Service {

	service := &Service{
		txManager:         txManager,
		conferenceService: conferenceService,
		worker:            worker,
		hasher:            hasher,
	}

	storage, ok := txManager.(repository.Storage)

	if ok {
		service.roomsRepo = storage.GetRoomsRepo()
		service.bookingRepo = storage.GetBookingRepo()
		service.authRepo = storage.GetAuthRepo()
	}
	return service
}
