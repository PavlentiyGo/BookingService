package postgres_repository

import (
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/repository"
	repository_errors "avitoBooking/internal/repository/erorrs"
	"avitoBooking/internal/repository/models"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueIdxOnBooking = "idx_bookings_one_active_per_slot"

type bookingRepository struct {
	TxManager repository.TxManager
}

func NewBookingRepository(
	txManager repository.TxManager,
) repository.BookingRepository {
	return &bookingRepository{TxManager: txManager}
}

func (r *bookingRepository) CreateBooking(
	ctx context.Context,
	booking domain.BookingRequest,
) (domain.BookingRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	INSERT INTO bookings(slot_id,user_id,conference_link)
	VALUES($1,$2,$3)
	RETURNING id,slot_id,user_id,status,conference_link,created_at;
	`
	var bookingModel models.BookingModel

	err := r.TxManager.GetExecutor(ctx).QueryRow(
		ctx,
		sqlQuery,
		booking.SlotId,
		booking.UserId,
		booking.ConferenceLink,
	).Scan(
		&bookingModel.Id,
		&bookingModel.SlotId,
		&bookingModel.UserId,
		&bookingModel.Status,
		&bookingModel.ConferenceLink,
		&bookingModel.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == uniqueIdxOnBooking {
					return domain.BookingRequest{}, repository_errors.ErrSlotIdConflict
				}
			}
			if pgErr.Code == "23503" {
				return domain.BookingRequest{}, repository_errors.ErrSlotNotFound
			}
		}
		return domain.BookingRequest{}, err
	}
	return models.BookingModelToDomain(bookingModel), nil
}
func (r *bookingRepository) GetSlotById(
	ctx context.Context,
	slotId uuid.UUID,
) (domain.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	SELECT id,room_id,start_time,end_time
	FROM slots 
	WHERE id = $1;
	`

	var model models.SlotModel

	err := r.TxManager.GetExecutor(ctx).QueryRow(
		ctx,
		sqlQuery,
		slotId).Scan(
		&model.Id,
		&model.RoomId,
		&model.StartTime,
		&model.EndTime,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Slot{}, fmt.Errorf("no rows in query: %w", repository_errors.ErrSlotNotFound)
		}
		return domain.Slot{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return models.SlotModelToDomain(model), nil
}

func (r *bookingRepository) ChangeBookingStatus(
	ctx context.Context,
	bookingId uuid.UUID,
	newStatus domain.BookingStatus,
) (domain.BookingRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	UPDATE bookings
	SET status = $1
	WHERE id = $2
	RETURNING id,slot_id,user_id,status,conference_link,created_at;
	`
	var newBooking models.BookingModel
	err := r.TxManager.GetExecutor(ctx).QueryRow(
		ctx,
		sqlQuery,
		newStatus,
		bookingId,
	).Scan(
		&newBooking.Id,
		&newBooking.SlotId,
		&newBooking.UserId,
		&newBooking.Status,
		&newBooking.ConferenceLink,
		&newBooking.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BookingRequest{}, repository_errors.ErrBookingNoExists
		}
		return domain.BookingRequest{}, err
	}
	return models.BookingModelToDomain(newBooking), nil
}
func (r *bookingRepository) GetBookingsByUserId(
	ctx context.Context,
	userId uuid.UUID,
) ([]domain.BookingRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	SELECT bookings.id,bookings.slot_id,bookings.user_id,bookings.status,bookings.conference_link,bookings.created_at
	FROM bookings
	JOIN slots
	ON slots.id = bookings.slot_id
	WHERE slots.start_time > now()
	AND bookings.user_id = $1;
	`

	rows, err := r.TxManager.GetExecutor(ctx).Query(ctx, sqlQuery, userId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.BookingRequest{}, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	var bookings []domain.BookingRequest

	for rows.Next() {
		var model models.BookingModel

		err = rows.Scan(
			&model.Id,
			&model.SlotId,
			&model.UserId,
			&model.Status,
			&model.ConferenceLink,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		bookings = append(bookings, models.BookingModelToDomain(model))
	}

	return bookings, nil
}

func (r *bookingRepository) GetBookingsCount(
	ctx context.Context,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	SELECT COUNT(*)
	FROM bookings;
	`
	var total int
	err := r.TxManager.GetExecutor(ctx).QueryRow(ctx, sqlQuery).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	return total, nil
}

func (r *bookingRepository) GetBookingsWithPagination(
	ctx context.Context,
	pagination domain.Pagination,
) ([]domain.BookingRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	offset := (pagination.Page - 1) * pagination.PageSize

	sqlQuery := `
	SELECT id,slot_id,user_id,status,conference_link,created_at
	FROM bookings
	ORDER BY created_at
	LIMIT $1 OFFSET $2;
	`

	rows, err := r.TxManager.GetExecutor(ctx).Query(
		ctx,
		sqlQuery,
		pagination.PageSize,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	var bookings []domain.BookingRequest

	for rows.Next() {
		var model models.BookingModel
		err = rows.Scan(
			&model.Id,
			&model.SlotId,
			&model.UserId,
			&model.Status,
			&model.ConferenceLink,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		bookings = append(bookings, models.BookingModelToDomain(model))
	}
	return bookings, nil
}
