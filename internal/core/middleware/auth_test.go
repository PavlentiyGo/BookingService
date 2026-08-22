package core_middleware

import (
	"avitoBooking/internal/core/auth"
	"avitoBooking/internal/core/auth/mocks"
	core_errors "avitoBooking/internal/core/errors"
	"avitoBooking/internal/core/http/response"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuth_InvalidRole(t *testing.T) {

	mockJwtProvider := mocks.NewJwtProvider(t)
	mockJwtProvider.On("ParseToken", "invalid-role").Return(auth.Token{
		MapClaims: nil,
		UserId:    uuid.Nil,
		Role:      "admin",
	}, nil)

	middleware := Auth(mockJwtProvider, "user")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	handler := middleware(nextHandler)

	headers := make(map[string]string, 1)
	headers["Authorization"] = "Bearer invalid-role"
	w := makeRequest(handler, "GET", "/auth", nil, headers)

	var resp response.ErrorResponse

	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusForbidden)
	assert.Contains(t, resp.Error.Message, core_errors.ErrForbidden.Error())
}
func TestAuth_NoToken(t *testing.T) {

	mockJwtProvider := mocks.NewJwtProvider(t)
	middleware := Auth(mockJwtProvider, "")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	handler := middleware(nextHandler)
	w := makeRequest(handler, "GET", "/auth", nil, nil)

	var resp response.ErrorResponse

	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	assert.Contains(t, resp.Error.Message, core_errors.ErrNotAuthorized.Error())
}
func TestAuth_InvalidTokenType(t *testing.T) {

	mockJwtProvider := mocks.NewJwtProvider(t)
	middleware := Auth(mockJwtProvider, "")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	handler := middleware(nextHandler)
	headers := make(map[string]string, 1)
	headers["Authorization"] = "invalid-type"
	w := makeRequest(handler, "GET", "/auth", nil, headers)

	var resp response.ErrorResponse

	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Contains(t, resp.Error.Message, core_errors.ErrWrongAuthType.Error())
}
func TestAuth_ParseTokenError(t *testing.T) {

	mockJwtProvider := mocks.NewJwtProvider(t)
	mockJwtProvider.On("ParseToken", "token").Return(auth.Token{}, errors.New("failed to parse token"))

	middleware := Auth(mockJwtProvider, "user")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	handler := middleware(nextHandler)

	headers := make(map[string]string, 1)
	headers["Authorization"] = "Bearer token"
	w := makeRequest(handler, "GET", "/auth", nil, headers)

	var resp response.ErrorResponse

	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	assert.NotEmpty(t, resp.Error.Message)
}
func TestAuth_ValidToken(t *testing.T) {

	mockJwtProvider := mocks.NewJwtProvider(t)
	mockJwtProvider.On("ParseToken", "token").Return(auth.Token{
		MapClaims: nil,
		Role:      "user",
		UserId:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
	}, nil)

	middleware := Auth(mockJwtProvider, "user")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := UserIdFromContext(r.Context())
		userRole := UserRoleFromContext(r.Context())
		if userId.String() != "00000000-0000-0000-0000-000000000002" || userRole != "user" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware(nextHandler)

	headers := make(map[string]string, 1)
	headers["Authorization"] = "Bearer token"
	w := makeRequest(handler, "GET", "/auth", nil, headers)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Empty(t, w.Body)
}
