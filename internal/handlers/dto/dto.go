package dto

import (
	"avitoBooking/internal/core/domain"
	"time"
)

type RoomDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	CreatedAt   string `json:"createdAt"`
}

type DummyLoginRequest struct {
	Role string `json:"role" validate:"required,oneof= user admin"`
}
type DummyLoginResponse struct {
	Token string `json:"token"`
}

type CreateRoomRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
	Capacity    *int    `json:"capacity" validate:"omitempty,numeric"`
}
type CreateRoomResponse struct {
	Room RoomDto `json:"room"`
}

type CreateScheduleRequest struct {
	DaysOfWeek []int  `json:"daysOfWeek" validate:"required,dive,numeric"`
	StartTime  string `json:"startTime" validate:"required"`
	EndTime    string `json:"endTime" validate:"required"`
}
type CreateScheduleResponse struct {
	ScheduleId string `json:"id"`
	RoomId     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type SlotDto struct {
	Id        string `json:"id"`
	RoomId    string `json:"roomId"`
	StartTime string `json:"startTime"`
	Endtime   string `json:"endTime"`
}
type CreateBookingRequest struct {
	SlotId         string `json:"slotId" validate:"required"`
	ConferenceLink bool   `json:"createConferenceLink"`
}
type BookingDto struct {
	Id             string               `json:"id"`
	SlotId         string               `json:"slotId"`
	UserId         string               `json:"userId"`
	Status         domain.BookingStatus `json:"status"`
	ConferenceLink string               `json:"conferenceLink"`
	CreatedAt      string               `json:"createdAt"`
}

type PaginationRequest struct {
	Page     *int `json:"page" validate:"omitempty,numeric"`
	PageSize *int `json:"pageSize" validate:"omitempty,numeric"`
}
type PaginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof= user admin"`
	Password string `json:"password" validate:"required"`
}
type UserDto struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
