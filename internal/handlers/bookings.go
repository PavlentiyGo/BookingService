package handlers

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	core_http_request "avitoBooking/internal/core/http/request"
	"avitoBooking/internal/core/http/response"
	core_middleware "avitoBooking/internal/core/middleware"
	"avitoBooking/internal/handlers/dto"
	"avitoBooking/internal/handlers/mapper"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (h *Handlers) CreateBooking(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	responser := response.NewResponser(w, ctx)

	var req dto.CreateBookingRequest
	if err := core_http_request.DecodeAndValidate(r, &req); err != nil {
		responser.ErrorResponse(err)
		return
	}
	slotId, err := uuid.Parse(req.SlotId)
	if err != nil {
		responser.ErrorResponse(fmt.Errorf("%s %w", err.Error(), core_errors.ErrInvalidSlotId))
		return
	}
	userId := core_middleware.UserIdFromContext(ctx)
	booking := domain.BookingRequest{
		SlotId:         slotId,
		UserId:         userId,
		ConferenceLink: nil,
	}
	if req.ConferenceLink {
		str := ""
		booking.ConferenceLink = &str
	}
	createdBooking, err := h.service.CreateBooking(ctx, booking)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	resp := mapper.BookingDomainToDto(createdBooking)
	jsonBody := map[string]dto.BookingDto{
		"booking": resp,
	}
	responser.WriteJson(http.StatusCreated, jsonBody)
}
func (h *Handlers) GetBookings(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	responser := response.NewResponser(w, ctx)
	var pageInt, pageSizeInt int
	var err error
	page := r.URL.Query().Get("page")
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			responser.ErrorResponse(fmt.Errorf("wrong page in query: %w: %w", err, core_errors.ErrInvalidRequest))
			return
		}
	}

	pageSize := r.URL.Query().Get("pageSize")
	if pageSize != "" {
		pageSizeInt, err = strconv.Atoi(pageSize)
		if err != nil {
			responser.ErrorResponse(fmt.Errorf("wrong pageSize in query: %w: %w", err, core_errors.ErrInvalidRequest))
			return
		}
	}
	pagination, err := domain.NewPagination(pageInt, pageSizeInt)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}

	bookings, bookingsTotal, err := h.service.GetBookings(ctx, pagination)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	bookingsDto := mapper.BookingsDomainToDto(bookings)
	paginationResp := dto.PaginationResponse{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		Total:    bookingsTotal,
	}
	resp := map[string]any{
		"bookings":   bookingsDto,
		"pagination": paginationResp,
	}
	responser.WriteJson(http.StatusOK, resp)

}
func (h *Handlers) GetMyBooking(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	responser := response.NewResponser(w, ctx)

	userId := core_middleware.UserIdFromContext(ctx)

	bookings, err := h.service.GetMyBookings(ctx, userId)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	dtos := mapper.BookingsDomainToDto(bookings)
	jsonBody := map[string][]dto.BookingDto{
		"bookings": dtos,
	}
	responser.WriteJson(http.StatusOK, jsonBody)
}
func (h *Handlers) CancelMyBooking(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	responser := response.NewResponser(w, ctx)

	bookingId := r.PathValue("bookingId")
	if bookingId == "" {
		responser.ErrorResponse(fmt.Errorf("invalid val in path: %w", core_errors.ErrInvalidBookingId))
		return
	}
	userId := core_middleware.UserIdFromContext(ctx)

	newBooking, err := h.service.CancelUserBooking(ctx, userId, bookingId)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	resp := mapper.BookingDomainToDto(newBooking)
	jsonBody := map[string]dto.BookingDto{
		"booking": resp,
	}
	responser.WriteJson(http.StatusOK, jsonBody)
}
