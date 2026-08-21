package models

import (
	"avitoBooking/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type RoomModel struct {
	Id          uuid.UUID
	Name        string
	Description *string
	Capacity    int
	CreatedAt   time.Time
}

type RoomScheduleModel struct {
	ScheduleId uuid.UUID
	RoomId     uuid.UUID
	DaysOfWeek []int
	StartTime  time.Time
	EndTime    time.Time
}

type SlotModel struct {
	Id        uuid.UUID
	RoomId    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}

type BookingModel struct {
	Id             uuid.UUID
	SlotId         uuid.UUID
	UserId         uuid.UUID
	Status         domain.BookingStatus
	ConferenceLink *string
	CreatedAt      time.Time
}
type UserModel struct {
	Id        uuid.UUID
	Role      string
	Email     string
	Password  []byte
	CreatedAt time.Time
}
