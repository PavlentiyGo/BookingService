package mapper

import (
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/handlers/dto"
)

func UserDomainToDto(
	user domain.User,
) dto.UserDto {
	return dto.UserDto{
		Id:        user.Id.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
