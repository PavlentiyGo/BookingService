package service

import (
	"avitoBooking/internal/core/domain"
	core_logger "avitoBooking/internal/core/logger"
	"avitoBooking/internal/repository"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const DefaultSlotDuration = time.Minute * 30

type Worker struct {
	TxManager        repository.TxManager
	roomsRepository  repository.RoomsRepository
	workerRepository repository.WorkerRepository
	SlotDuration     time.Duration
	WorkerCooldown   time.Duration
}

type WorkerInterface interface {
	Start(ctx context.Context)
	CreateSlots(ctx context.Context, roomSchedule *domain.RoomSchedule) error
}

func NewWorker(txManager repository.TxManager) *Worker {

	worker := Worker{
		TxManager:      txManager,
		SlotDuration:   DefaultSlotDuration,
		WorkerCooldown: time.Hour * 24, // TODO FROM Config
	}

	if tx, ok := txManager.(repository.Storage); ok {
		worker.roomsRepository = tx.GetRoomsRepo()
		worker.workerRepository = tx.GetWorkerRepo()
	}

	return &worker

}

func (w *Worker) Start(ctx context.Context) {
	logger := core_logger.FromContext(ctx)
	Log := logger.With(zap.String("func", "worker_start"))
	Log.Debug("starting")
	err := w.workerRepository.CleanSlots(ctx)
	if err != nil {
		Log.Error("failed to clean slots", zap.Error(err))
	}
	err = w.CreateSlots(ctx, nil)
	if err != nil {
		Log.Error("failed to create new slots", zap.Error(err))
	}
	Log.Debug("ending")

	timer := time.NewTimer(w.WorkerCooldown)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			Log.Debug("closing worker")
			return
		case <-timer.C:
			Log.Debug("starting work")

			err = w.workerRepository.CleanSlots(ctx)
			if err != nil {
				Log.Error("failed to clean slots", zap.Error(err))
			}
			err = w.CreateSlots(ctx, nil)
			if err != nil {
				Log.Error("failed to create new slots", zap.Error(err))
			}
			Log.Debug("end working")
		}
	}
}
func (w *Worker) CreateSlots(ctx context.Context, roomSchedule *domain.RoomSchedule) error {
	var schedules []domain.RoomSchedule
	var err error
	if roomSchedule == nil {
		schedules, err = w.roomsRepository.GetRoomsSchedule(ctx) // TODO добавить ручку на получение расписания комнат которым нужны слоты
	} else {
		schedules = append(schedules, *roomSchedule)
	}
	if err != nil {
		return fmt.Errorf("failed to get rooms: %w", err)
	}
	var slots []domain.Slot
	for i := 0; i < len(schedules); i++ {
		var slot domain.Slot
		startTime := schedules[i].StartTime
		endtime := schedules[i].EndTime
		slot.RoomId = schedules[i].RoomId
		var times [][2]time.Time
		for !startTime.Add(w.SlotDuration).After(endtime) {
			slotTimeStart := time.Date(0, 0, 0, startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
			startTime = startTime.Add(w.SlotDuration)
			slotTimeEnd := time.Date(0, 0, 0, startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
			times = append(times, [2]time.Time{slotTimeStart, slotTimeEnd})
		}

		for j := 0; j < len(schedules[i].DaysOfWeek); j++ {
			day := w.calcNextWeekday(schedules[i].DaysOfWeek[j])
			for _, val := range times {
				slot.StartTime = time.Date(day.Year(), day.Month(), day.Day(), val[0].Hour(), val[0].Minute(), 0, 0, time.UTC)
				slot.EndTime = time.Date(day.Year(), day.Month(), day.Day(), val[1].Hour(), val[1].Minute(), 0, 0, time.UTC)
				slots = append(slots, slot)
			}
		}
	}
	err = w.TxManager.WithinTx(ctx, func(txCtx context.Context) error {
		err = w.workerRepository.CreateSlots(txCtx, slots)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (w *Worker) calcNextWeekday(weekdayNumber int) time.Time {

	now := time.Now().Weekday()
	if now == 0 {
		now = 7
	}
	var difference int
	if int(now) <= weekdayNumber {
		difference = weekdayNumber - int(now)
	} else {
		difference = 7 - (int(now) - weekdayNumber)
	}

	day := time.Now().AddDate(0, 0, difference)

	nextDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return nextDay
}

//# 1. Получаем токен один раз и сохраняем в переменную
//TOKEN=$(curl -s -X POST http://localhost:8080/dummyLogin -d '{"role": "user"}' | jq -r '.token')
//
//# 2. Запускаем vegeta, подставляя токен в хедер
//echo "GET http://localhost:8080/api/rooms/00000000-0000-0000-0000-000000000001/slots?date=2026-08-20" | \
//vegeta attack -duration=30s -rate=100/1s -header="Authorization: Bearer $TOKEN" | \
//vegeta report -type=text
