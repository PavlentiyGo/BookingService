package handlers

import (
	"avitoBooking/internal/core/auth"
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/handlers/dto"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateBooking_Success(t *testing.T) {
	deps := newTestDeps(t)
	slotId := uuid.New()
	userId := userDummyUUID

	token := "valid-user-token"
	deps.jwtProvider.On("ParseToken", token).Return(auth.Token{
		UserId: userId,
		Role:   "user",
	}, nil)

	reqDto := dto.CreateBookingRequest{
		SlotId:         slotId.String(),
		ConferenceLink: true,
	}
	body, _ := json.Marshal(reqDto)

	expectedBooking := domain.BookingRequest{
		Id:     uuid.New(),
		SlotId: slotId,
		UserId: userId,
		Status: domain.BookingStatusActive,
	}

	deps.service.On("CreateBooking", mock.Anything, mock.MatchedBy(func(b domain.BookingRequest) bool {
		return b.SlotId == slotId && b.UserId == userId
	})).Return(expectedBooking, nil)

	headers := map[string]string{"Authorization": "Bearer " + token}
	w := makeRequest(deps.router, "POST", "/bookings/create", body, headers)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]dto.BookingDto
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, expectedBooking.Id.String(), resp["booking"].Id)
}

func TestGetMyBookings_Success(t *testing.T) {
	deps := newTestDeps(t)
	userId := userDummyUUID
	token := "user-token"

	deps.jwtProvider.On("ParseToken", token).Return(auth.Token{
		UserId: userId,
		Role:   "user",
	}, nil)

	mockBookings := []domain.BookingRequest{
		{Id: uuid.New(), UserId: userId, Status: domain.BookingStatusActive},
	}

	deps.service.On("GetMyBookings", mock.Anything, userId).Return(mockBookings, nil)

	headers := map[string]string{"Authorization": "Bearer " + token}
	w := makeRequest(deps.router, "GET", "/bookings/my", nil, headers)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]dto.BookingDto
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Len(t, resp["bookings"], 1)
	assert.Equal(t, mockBookings[0].Id.String(), resp["bookings"][0].Id)
}
