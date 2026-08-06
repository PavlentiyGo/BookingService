package mapper

import (
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/handlers/dto"
)

func RoomDtoToDomain(room dto.CreateRoomRequest) domain.Room {
	return domain.Room{
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
	}
}
func RoomsDomainToDto(rooms []domain.Room) []dto.RoomDto {
	dtoRooms := make([]dto.RoomDto, len(rooms))
	for i := 0; i < len(rooms); i++ {
		dtoRooms[i] = RoomDomainToDto(rooms[i])
	}
	return dtoRooms
}

func RoomDomainToDto(room domain.Room) dto.RoomDto {
	roomDto := dto.RoomDto{
		Id:        room.Id.String(),
		Name:      room.Name,
		CreatedAt: room.CreatedAt.String(),
	}
	if room.Description != nil {
		roomDto.Description = *room.Description
	}
	if room.Capacity != nil {
		roomDto.Capacity = *room.Capacity
	}
	return roomDto
}
