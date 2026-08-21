package handlers

import (
	"avitoBooking/internal/core/domain"
	core_http_request "avitoBooking/internal/core/http/request"
	"avitoBooking/internal/core/http/response"
	"avitoBooking/internal/handlers/dto"
	"avitoBooking/internal/handlers/mapper"
	"net/http"

	"github.com/google/uuid"
)

var (
	adminDummyUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userDummyUUID  = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

func (h *Handlers) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	responser := response.NewResponser(w, ctx)

	var req dto.RegisterRequest
	if err := core_http_request.DecodeAndValidate(r, &req); err != nil {
		responser.ErrorResponse(err)
		return
	}
	user := domain.User{
		Role:     req.Role,
		Email:    req.Email,
		Password: []byte(req.Password),
	}
	registeredUser, err := h.service.Register(ctx, user)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}
	resp := map[string]dto.UserDto{
		"user": mapper.UserDomainToDto(registeredUser),
	}
	responser.WriteJson(http.StatusCreated, resp)
}

func (h *Handlers) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	//ctx := r.Context()
	//responser := response.NewResponser(w, ctx)

}

func (h *Handlers) DummyLogin(
	w http.ResponseWriter,
	r *http.Request,
) {

	responser := response.NewResponser(w, r.Context())

	request := dto.DummyLoginRequest{}

	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responser.ErrorResponse(err)
		return
	}

	var uuid uuid.UUID
	switch request.Role {
	case "user":
		uuid = userDummyUUID
	case "admin":
		uuid = adminDummyUUID
	}

	token, err := h.jwtProvider.NewToken(request.Role, uuid)
	if err != nil {
		responser.ErrorResponse(err)
		return
	}

	resp := dto.DummyLoginResponse{Token: token}

	responser.WriteJson(http.StatusOK, resp)
}
