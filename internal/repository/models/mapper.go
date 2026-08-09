package models

import (
	"avitoBooking/internal/core/domain"
)

func RoomModelToDomain(
	room RoomModel,
) domain.Room {
	return domain.Room{
		Id:          room.Id,
		Name:        room.Name,
		Description: room.Description,
		Capacity:    &room.Capacity,
		CreatedAt:   &room.CreatedAt,
	}
}

func RoomModelsToDomain(
	rooms []RoomModel,
) []domain.Room {
	domainRooms := make([]domain.Room, len(rooms))
	for i := 0; i < len(rooms); i++ {
		domainRooms[i] = RoomModelToDomain(rooms[i])
	}
	return domainRooms
}

func SlotModelToDomain(
	slot SlotModel,
) domain.Slot {
	return domain.Slot{
		Id:        slot.Id,
		RoomId:    slot.RoomId,
		StartTime: slot.StartTime,
		EndTime:   slot.EndTime,
	}
}
func BookingModelToDomain(
	booking BookingModel,
) domain.BookingRequest {
	return domain.BookingRequest{
		Id:             booking.Id,
		SlotId:         booking.SlotId,
		UserId:         booking.UserId,
		Status:         booking.Status,
		ConferenceLink: booking.ConferenceLink,
		CreatedAt:      booking.CreatedAt,
	}

}
