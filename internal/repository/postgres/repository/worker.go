package postgres_repository

import (
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/repository"
	"context"
	"fmt"
	"strings"
	"time"
)

type workerRepository struct {
	TxManager repository.TxManager
}

func NewWorkerRepository(
	txManager repository.TxManager,
) repository.WorkerRepository {
	return &workerRepository{TxManager: txManager}
}

func (r *workerRepository) CreateSlots(
	ctx context.Context,
	slots []domain.Slot,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := r.buildSlotsQuery(slots)
	_, err := r.TxManager.GetExecutor(ctx).Exec(ctx, sqlQuery)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
func (r *workerRepository) buildSlotsQuery(slots []domain.Slot) string {
	var builder strings.Builder
	sqlQuery := `
INSERT INTO slots(room_id,start_time,end_time)
VALUES
	`

	builder.WriteString(sqlQuery)
	for ind, slot := range slots {
		if ind != 0 {
			builder.WriteString(",\n")
		}
		builder.WriteString(fmt.Sprintf("('%s','%s','%s')", slot.RoomId.String(), slot.StartTime.Format(time.RFC3339), slot.EndTime.Format(time.RFC3339))) // TODO исправить вместо .Format
	}
	endQuery := `ON CONFLICT DO NOTHING;`
	builder.WriteString("\n")
	builder.WriteString(endQuery)

	return builder.String()
}
func (r *workerRepository) CleanSlots(
	ctx context.Context,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	DELETE FROM slots
	WHERE end_time < now();	
	`

	_, err := r.TxManager.GetExecutor(ctx).Exec(ctx, sqlQuery)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
