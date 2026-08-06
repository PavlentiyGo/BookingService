package handlers

import (
	core_http_request "avitoBooking/internal/core/http/request"
	"avitoBooking/internal/core/http/response"
	"avitoBooking/internal/handlers/dto"
	"net/http"
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

}
func (h *Handlers) GetBookings(
	w http.ResponseWriter,
	r *http.Request,
) {

}
func (h *Handlers) GetMyBooking(
	w http.ResponseWriter,
	r *http.Request,
) {

}
func (h *Handlers) CancelMyBooking(
	w http.ResponseWriter,
	r *http.Request,
) {

}
