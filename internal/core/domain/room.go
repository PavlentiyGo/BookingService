package domain

import (
	core_errors "avitoBooking/internal/core/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id          uuid.UUID
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   *time.Time
}

func (r *Room) Validate() error {
	if r.Capacity != nil {
		if *r.Capacity < 0 {
			return fmt.Errorf("wrong room capacity: %w", core_errors.ErrInvalidRoomCapacity)
		}
	}
	return nil
}

type RoomSchedule struct {
	ScheduleId *uuid.UUID
	RoomId     uuid.UUID
	DaysOfWeek []int
	StartTime  time.Time
	EndTime    time.Time
}

func (s *RoomSchedule) Validate() error {
	if s.RoomId == uuid.Nil {
		return core_errors.ErrWrongRoomId
	}
	if len(s.DaysOfWeek) > 7 {
		return core_errors.ErrInvalidDays
	}
	for i := 0; i < len(s.DaysOfWeek); i++ {
		if s.DaysOfWeek[i] <= 0 || s.DaysOfWeek[i] > 7 {
			return core_errors.ErrInvalidDays
		}
	}
	if s.StartTime.After(s.EndTime) {
		return core_errors.ErrInvalidTime
	}
	return nil
}

type Slot struct {
	Id        uuid.UUID
	RoomId    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}
