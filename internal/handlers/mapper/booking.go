package mapper

import (
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/handlers/dto"
)

func BookingDomainToDto(
	booking domain.BookingRequest,
) dto.BookingDto {
	bookingResp := dto.BookingDto{
		Id:        booking.Id.String(),
		SlotId:    booking.SlotId.String(),
		UserId:    booking.UserId.String(),
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt.String(),
	}
	if booking.ConferenceLink != nil {
		bookingResp.ConferenceLink = *booking.ConferenceLink
	}
	return bookingResp
}
func BookingsDomainToDto(
	bookings []domain.BookingRequest,
) []dto.BookingDto {

	bookingsDtos := make([]dto.BookingDto, len(bookings))
	for i := 0; i < len(bookings); i++ {
		bookingsDtos[i] = BookingDomainToDto(bookings[i])
	}
	return bookingsDtos
}
