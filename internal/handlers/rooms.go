package handlers

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	core_http_request "avitoBooking/internal/core/http/request"
	"avitoBooking/internal/core/http/response"
	"avitoBooking/internal/handlers/dto"
	"avitoBooking/internal/handlers/mapper"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *Handlers) GetRooms(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	responser := response.NewResponser(w, ctx)

	rooms, err := h.service.GetRooms(ctx)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	dtos := mapper.RoomsDomainToDto(rooms)
	jsonBody := map[string][]dto.RoomDto{
		"rooms": dtos,
	}
	responser.WriteJson(http.StatusOK, jsonBody)
}

func (h *Handlers) CreateRoom(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	responser := response.NewResponser(w, ctx)

	req := dto.CreateRoomRequest{}
	if err := core_http_request.DecodeAndValidate(r, &req); err != nil {
		responser.ErrorResponse(err)
		return
	}

	createdRoom, err := h.service.CreateRoom(
		ctx,
		mapper.RoomDtoToDomain(req),
	)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	dtoResp := mapper.RoomDomainToDto(createdRoom)

	jsonBody := map[string]dto.RoomDto{
		"room": dtoResp,
	}
	responser.WriteJson(http.StatusCreated, jsonBody)
}
func (h *Handlers) CreateSchedule(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	responser := response.NewResponser(w, ctx)

	roomId := r.PathValue("roomId")
	roomUUID, err := uuid.Parse(roomId)
	if err != nil {
		responser.ErrorResponse(fmt.Errorf("%w: %w", err, core_errors.ErrWrongRoomId))
		return
	}
	var req dto.CreateScheduleRequest
	if err = core_http_request.DecodeAndValidate(r, &req); err != nil {
		responser.ErrorResponse(err)
		return
	}
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		responser.ErrorResponse(fmt.Errorf("%w: %w", err, core_errors.ErrInvalidTime))
		return
	}
	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		responser.ErrorResponse(fmt.Errorf("%w: %w", err, core_errors.ErrInvalidTime))
		return
	}
	roomScheduleDomain := domain.RoomSchedule{
		RoomId:     roomUUID,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}
	fmt.Println(roomScheduleDomain.StartTime, roomScheduleDomain.EndTime)
	createdSchedule, err := h.service.CreateSchedule(ctx, roomScheduleDomain)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	resp := map[string]dto.CreateScheduleResponse{
		"schedule": {
			ScheduleId: createdSchedule.ScheduleId.String(),
			RoomId:     createdSchedule.RoomId.String(),
			DaysOfWeek: createdSchedule.DaysOfWeek,
			StartTime:  createdSchedule.StartTime.Format("15:04"),
			EndTime:    createdSchedule.EndTime.Format("15:04"),
		},
	}

	responser.WriteJson(http.StatusOK, resp)
}

func (h *Handlers) GetSlots(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	responser := response.NewResponser(w, ctx)

	pathVal := r.PathValue("roomId")
	roomId, err := uuid.Parse(pathVal)
	if err != nil {
		responser.ErrorResponse(fmt.Errorf("%w: %w", err, core_errors.ErrWrongRoomId))
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		responser.ErrorResponse(core_errors.ErrMissingDate)
		return
	}
	slots, err := h.service.GetSlots(ctx, roomId, date)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	slotsDto := mapper.SlotsDomainToDto(slots)
	jsonBody := map[string][]dto.SlotDto{
		"slots": slotsDto,
	}

	responser.WriteJson(http.StatusOK, jsonBody)

}
