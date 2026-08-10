package service

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	core_logger "avitoBooking/internal/core/logger"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) CreateRoom(
	ctx context.Context,
	room domain.Room,
) (domain.Room, error) {
	if err := room.Validate(); err != nil {
		return domain.Room{}, err
	}
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
	Log := core_logger.FromContext(ctx)
	Log = Log.With(zap.String("func", "service-CreateSchedule"))

	if err := schedule.Validate(); err != nil {
		return domain.RoomSchedule{}, err
	}
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
	var createdSchedule domain.RoomSchedule
	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		roomSchedule, err := s.roomsRepo.CreateSchedule(ctx, schedule)
		if err != nil {
			return fmt.Errorf("failed to create room schedule: %w", err)
		}
		createdSchedule = roomSchedule
		err = s.worker.CreateSlots(txCtx, &roomSchedule)
		if err != nil {
			Log.Error("failed to create slots for room schedule", zap.Error(err))
		}
		return nil
	})
	if err != nil {
		return domain.RoomSchedule{}, fmt.Errorf("error during tx: %w", err)
	}
	return createdSchedule, nil
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
	if dateTime.Before(time.Now()) {
		return nil, fmt.Errorf("date time must be after current date: %w", core_errors.ErrInvalidDateTime)
	}

	return s.roomsRepo.GetSlots(ctx, roomId, dateTime)
}
