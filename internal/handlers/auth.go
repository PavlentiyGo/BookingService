package handlers

import (
	"avitoBooking/internal/core/http/response"
	"avitoBooking/internal/handlers/dto"
	"encoding/json"
	"net/http"
)

func (h *Handlers) Register(
	w http.ResponseWriter,
	r *http.Request,
) {

}

func (h *Handlers) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

}

func (h *Handlers) DummyLogin(
	w http.ResponseWriter,
	r *http.Request,
) {

	responser := response.NewResponser(w, r.Context())

	request := dto.DummyLoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responser.ErrorResponse()
	}

}
