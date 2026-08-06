package service

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Service) CreateRoom(
	ctx context.Context,
	room domain.Room,
) (domain.Room, error) {
	if err := room.Validate(); err != nil {
		return domain.Room{}, err
	}
	newId, err := uuid.NewUUID()
	if err != nil {
		return domain.Room{}, fmt.Errorf("failed to create uuid for room")
	}
	room.Id = newId
	newRoom, err := s.roomsRepo.CreateRoom(ctx, room)
	if err != nil {
		return domain.Room{}, err
	}
	return newRoom, nil
}
func (s *Service) GetRooms(
	ctx context.Context,
) ([]domain.Room, error) {
	return s.roomsRepo.GetRooms(ctx)
}

func (s *Service) CreateSchedule(
	ctx context.Context,
	schedule domain.RoomSchedule,
) (domain.RoomSchedule, error) {
	if err := schedule.Validate(); err != nil {
		return domain.RoomSchedule{}, err
	}
	newUuid, err := uuid.NewUUID()
	if err != nil {
		return domain.RoomSchedule{}, fmt.Errorf("failed to create new uuid for schedule: %w", err)
	}
	schedule.ScheduleId = &newUuid
	now := time.Now()
	schedule.StartTime = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		schedule.StartTime.Hour(),
		schedule.StartTime.Minute(),
		0, 0, time.UTC,
	)
	schedule.EndTime = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		schedule.EndTime.Hour(),
		schedule.EndTime.Minute(),
		0, 0, time.UTC,
	)
	return s.roomsRepo.CreateSchedule(ctx, schedule)
}
func (s *Service) GetSlots(
	ctx context.Context,
	roomId uuid.UUID,
	date string,
) ([]domain.Slot, error) {

	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dateTime: %w: %w", err, core_errors.ErrInvalidDateTime)
	}
	return s.roomsRepo.GetSlots(ctx, roomId, dateTime)
}
