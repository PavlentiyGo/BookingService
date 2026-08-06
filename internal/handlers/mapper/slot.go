package mapper

import (
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/handlers/dto"
)

func SlotDomainToDto(
	slot domain.Slot,
) dto.SlotDto {
	return dto.SlotDto{
		Id:        slot.Id.String(),
		RoomId:    slot.RoomId.String(),
		StartTime: slot.StartTime.String(),
		Endtime:   slot.EndTime.String(),
	}
}
func SlotsDomainToDto(
	slots []domain.Slot,
) []dto.SlotDto {
	dtoSlots := make([]dto.SlotDto, len(slots))
	for i := 0; i < len(slots); i++ {
		dtoSlots[i] = SlotDomainToDto(slots[i])
	}

	return dtoSlots
}
