package handlers

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	"avitoBooking/internal/handlers/dto"
	"avitoBooking/internal/handlers/mapper"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDummyLogin_InvalidRole(t *testing.T) {
	deps := newTestDeps(t)
	body := dto.DummyLoginRequest{
		Role: "seller",
	}
	bytes, err := json.Marshal(body)
	require.NoError(t, err)
	w := makeRequest(deps.router, "POST", "/dummyLogin", bytes, nil)
	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestDummyLogin_ValidRoles(t *testing.T) {
	deps := newTestDeps(t)
	cases := []dto.DummyLoginRequest{
		{
			Role: "user",
		},
		{
			Role: "admin",
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.Role, func(t *testing.T) {
			bytes, err := json.Marshal(tCase)
			require.NoError(t, err)
			deps.jwtProvider.On(
				"NewToken",
				tCase.Role,
				mock.Anything,
			).Return(
				"valid_token",
				nil,
			)
			w := makeRequest(deps.router, "POST", "/dummyLogin", bytes, nil)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "valid_token")
		})
	}
}
func TestDummyLogin_InternalError(t *testing.T) {
	deps := newTestDeps(t)
	body := dto.DummyLoginRequest{
		Role: "user",
	}
	bytes, err := json.Marshal(body)
	require.NoError(t, err)
	deps.jwtProvider.On(
		"NewToken",
		body.Role,
		mock.Anything,
	).Return(
		"",
		errors.New("internal_error"),
	)
	w := makeRequest(deps.router, "POST", "/dummyLogin", bytes, nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), core_errors.ErrInternalError.Error())
}

func TestRegister_Success(t *testing.T) {
	deps := newTestDeps(t)

	dtoReq := dto.RegisterRequest{
		Email:    "pavel@gmail.com",
		Role:     "user",
		Password: "password",
	}
	body, err := json.Marshal(dtoReq)
	require.NoError(t, err)
	domainUser := domain.User{
		Email:    dtoReq.Email,
		Role:     dtoReq.Role,
		Password: []byte(dtoReq.Password),
	}
	registeredUser := domain.User{
		Id:        uuid.Nil,
		Role:      domainUser.Role,
		Email:     domainUser.Email,
		Password:  domainUser.Password,
		CreatedAt: time.Time{},
	}
	deps.service.On(
		"Register",
		mock.Anything,
		domainUser,
	).Return(
		registeredUser,
		nil,
	)
	w := makeRequest(deps.router, "POST", "/register", body, nil)

	mustResp := map[string]dto.UserDto{
		"user": mapper.UserDomainToDto(registeredUser),
	}
	resp := make(map[string]dto.UserDto)
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusCreated)
	assert.Equal(t, resp, mustResp)
}
func TestRegister_InvalidBody(t *testing.T) {
	deps := newTestDeps(t)
	cases := map[string]dto.RegisterRequest{
		"invalid_email": {
			Email:    "pavel",
			Role:     "user",
			Password: "password",
		},
		"invalid_role": {
			Email:    "pavel@gmail.com",
			Role:     "seller",
			Password: "password",
		},
		"invalid_password": {
			Email:    "pavel@gmail.com",
			Role:     "user",
			Password: "",
		},
	}
	for name, tCase := range cases {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(tCase)
			require.NoError(t, err)
			w := makeRequest(deps.router, "POST", "/register", body, nil)
			assert.Equal(t, w.Code, http.StatusBadRequest)
			assert.NotEmpty(t, w.Body.String())
		})
	}
}
func TestRegister_ServiceError(t *testing.T) { // TODO add more cases
	deps := newTestDeps(t)

	dtoReq := dto.RegisterRequest{
		Email:    "pavel@gmail.com",
		Role:     "user",
		Password: "password",
	}
	body, err := json.Marshal(dtoReq)
	require.NoError(t, err)
	domainUser := domain.User{
		Email:    dtoReq.Email,
		Role:     dtoReq.Role,
		Password: []byte(dtoReq.Password),
	}
	deps.service.On(
		"Register",
		mock.Anything,
		domainUser,
	).Return(
		domain.User{},
		core_errors.ErrEmailExists,
	)
	w := makeRequest(deps.router, "POST", "/register", body, nil)

	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), core_errors.ErrEmailExists.Error())
}
