package postgres_repository

import (
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/repository"
	repository_errors "avitoBooking/internal/repository/erorrs"
	"avitoBooking/internal/repository/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type roomsRepository struct {
	TxManager repository.TxManager
}

func NewRoomsRepository(
	txManager repository.TxManager,
) repository.RoomsRepository {
	return &roomsRepository{TxManager: txManager}
}

func (r *roomsRepository) CreateRoom(
	ctx context.Context,
	room domain.Room,
) (domain.Room, error) {

	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	INSERT INTO rooms(name,description,capacity)
	VALUES($1,$2,$3)
	RETURNING id,name,description,capacity,createdAt;
	`
	row := r.TxManager.GetExecutor(ctx).QueryRow(
		ctx, sqlQuery,
		room.Name,
		room.Description,
		room.Capacity,
	)

	var model models.RoomModel
	err := row.Scan(
		&model.Id,
		&model.Name,
		&model.Description,
		&model.Capacity,
		&model.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.Room{}, repository_errors.ErrRoomAlreadyExists
			}
		}
		return domain.Room{}, fmt.Errorf("row scan: %w", err)
	}
	createdRoom := domain.Room{
		Id:          model.Id,
		Name:        model.Name,
		Description: model.Description,
		Capacity:    &model.Capacity,
		CreatedAt:   &model.CreatedAt,
	}
	return createdRoom, nil
}
func (r *roomsRepository) GetRooms(
	ctx context.Context,
) ([]domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	SELECT id,name,description,capacity,createdAt 
	FROM rooms;
	`

	var roomModels []models.RoomModel

	rows, err := r.TxManager.GetExecutor(ctx).Query(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}
	for rows.Next() {
		var room models.RoomModel

		err = rows.Scan(
			&room.Id,
			&room.Name,
			&room.Description,
			&room.Capacity,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		roomModels = append(roomModels, room)
	}
	return models.RoomModelsToDomain(roomModels), nil
}
func (r *roomsRepository) CreateSchedule(
	ctx context.Context,
	roomSchedule domain.RoomSchedule,
) (domain.RoomSchedule, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	INSERT INTO schedules(room_id,days_of_week,start_time,end_time)
	VALUES($1,$2,$3,$4)
	RETURNING id,room_id,days_of_week,start_time,end_time;
	`
	var model models.RoomScheduleModel
	err := r.TxManager.GetExecutor(ctx).QueryRow(
		ctx,
		sqlQuery,
		roomSchedule.RoomId,
		roomSchedule.DaysOfWeek,
		roomSchedule.StartTime,
		roomSchedule.EndTime,
	).Scan(
		&model.ScheduleId,
		&model.RoomId,
		&model.DaysOfWeek,
		&model.StartTime,
		&model.EndTime,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.RoomSchedule{}, repository_errors.ErrRoomScheduleExists
			}
			if pgErr.Code == "23503" {
				return domain.RoomSchedule{}, repository_errors.ErrRoomNotFound
			}
		}
		return domain.RoomSchedule{}, fmt.Errorf("row scan: %w", err)
	}
	newRoomSchedule := domain.RoomSchedule{
		ScheduleId: &model.ScheduleId,
		RoomId:     model.RoomId,
		DaysOfWeek: model.DaysOfWeek,
		StartTime:  model.StartTime,
		EndTime:    model.EndTime,
	}

	return newRoomSchedule, nil
}

func (r *roomsRepository) GetSlots(
	ctx context.Context,
	roomId uuid.UUID,
	dateTime time.Time,
) ([]domain.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	SELECT slots.id,slots.room_id,slots.start_time,slots.end_time FROM slots
	LEFT JOIN bookings
	ON slots.id = bookings.slot_id
	WHERE 'cancelled' = COALESCE(bookings.status,'cancelled')
	AND slots.room_id = $1
	AND slots.start_time::date = $2::date;
	`
	rows, err := r.TxManager.GetExecutor(ctx).Query(
		ctx,
		sqlQuery,
		roomId,
		dateTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}
	var slotsDomain []domain.Slot

	for rows.Next() {
		var model models.SlotModel
		err = rows.Scan(
			&model.Id,
			&model.RoomId,
			&model.StartTime,
			&model.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		slotsDomain = append(slotsDomain, models.SlotModelToDomain(model))
	}
	return slotsDomain, nil
}

func (r *roomsRepository) GetRoomsSchedule(
	ctx context.Context,
) ([]domain.RoomSchedule, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	SELECT id,room_id,days_of_week,start_time,end_time
	FROM schedules;
	`

	rows, err := r.TxManager.GetExecutor(ctx).Query(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	var schedules []domain.RoomSchedule
	for rows.Next() {
		var model models.RoomScheduleModel

		err = rows.Scan(
			&model.ScheduleId,
			&model.RoomId,
			&model.DaysOfWeek,
			&model.StartTime,
			&model.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		schedules = append(schedules, models.RoomScheduleModelToDomain(model))
	}
	return schedules, nil
}
