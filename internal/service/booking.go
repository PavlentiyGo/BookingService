package service

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	core_logger "avitoBooking/internal/core/logger"
	"time"

	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) CreateBooking(
	ctx context.Context,
	booking domain.BookingRequest,
) (domain.BookingRequest, error) {
	Log := core_logger.FromContext(ctx)
	Log = Log.With(zap.String("func", "createBooking"))
	link, err := s.conferenceService.CreateConferenceLink(ctx, booking.SlotId)

	if err != nil {
		return domain.BookingRequest{}, fmt.Errorf("failed to create conference link: %w", err)
	}
	var answ domain.BookingRequest
	booking.ConferenceLink = &link
	err = s.txManager.WithinTx(ctx, func(txCtx context.Context) error {

		slot, err := s.bookingRepo.GetSlotById(ctx, booking.SlotId)
		if err != nil {
			return fmt.Errorf("failed to get slot by id: %w", err)
		}
		if slot.StartTime.Before(time.Now()) {
			return core_errors.ErrSlotInThePast
		}
		newBooking, err := s.bookingRepo.CreateBooking(ctx, booking)
		if err != nil {
			newErr := s.conferenceService.CancelConference(nil, "")
			if newErr != nil {
				Log.Error("failed to cancel conference link", zap.String("link", link), zap.Error(newErr))
			}
			return fmt.Errorf("failed to create booking: %w", err)
		}
		answ = newBooking
		return nil
	})
	if err != nil {
		return domain.BookingRequest{}, err
	}
	return answ, nil
}
func (s *Service) GetMyBookings(
	ctx context.Context,
	userId uuid.UUID,
) ([]domain.BookingRequest, error) {
	return s.bookingRepo.GetBookingsByUserId(ctx, userId)
}
func (s *Service) CancelUserBooking(
	ctx context.Context,
	userId uuid.UUID,
	bookingId string,
) (domain.BookingRequest, error) {

	bookingUUID, err := uuid.Parse(bookingId)
	if err != nil {
		return domain.BookingRequest{}, fmt.Errorf("failed to parser booking id: %w", core_errors.ErrInvalidBookingId)
	}

	var booking domain.BookingRequest

	err = s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		newBooking, err := s.bookingRepo.ChangeBookingStatus(ctx, bookingUUID, domain.BookingStatusCancelled)
		if err != nil {
			return fmt.Errorf("failed to change booking status: %w", err)
		}
		if newBooking.UserId != userId {
			return fmt.Errorf("this is not your bookings: %w", core_errors.ErrForbidden)
		}
		booking = newBooking
		return nil
	})
	if err != nil {
		return domain.BookingRequest{}, fmt.Errorf("error during tx: %w", err)
	}
	return booking, nil
}
func (s *Service) GetBookings(
	ctx context.Context,
	pagination domain.Pagination,
) ([]domain.BookingRequest, int, error) {

	var bookingsCount int
	var bookings []domain.BookingRequest
	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		total, err := s.bookingRepo.GetBookingsCount(ctx)
		if err != nil {
			return fmt.Errorf("failed to get all bookings: %w", err)
		}
		bookingsCount = total
		gotBookings, err := s.bookingRepo.GetBookingsWithPagination(ctx, pagination)
		if err != nil {
			return fmt.Errorf("failed to get all bookings: %w", err)
		}
		bookings = gotBookings
		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("error during tx: %w", err)
	}
	return bookings, bookingsCount, nil
}
