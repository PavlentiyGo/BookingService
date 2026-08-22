package handlers

import (
	"avitoBooking/internal/core/auth"
	"avitoBooking/internal/core/domain"
	"avitoBooking/internal/handlers/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateRoom_AdminSuccess(t *testing.T) {
	deps := newTestDeps(t)
	adminId := adminDummyUUID
	token := "admin-token"

	deps.jwtProvider.On("ParseToken", token).Return(auth.Token{
		UserId: adminId,
		Role:   "admin",
	}, nil)

	capacity := 10
	desc := "Conference room"
	reqDto := dto.CreateRoomRequest{
		Name:        "Room 1",
		Capacity:    &capacity,
		Description: &desc,
	}
	body, _ := json.Marshal(reqDto)

	now := time.Now()
	createdRoom := domain.Room{
		Id:          uuid.New(),
		Name:        reqDto.Name,
		Capacity:    &capacity,
		Description: &desc,
		CreatedAt:   &now,
	}

	deps.service.On("CreateRoom", mock.Anything, mock.Anything).Return(createdRoom, nil)

	headers := map[string]string{"Authorization": "Bearer " + token}
	w := makeRequest(deps.router, "POST", "/rooms/create", body, headers)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]dto.RoomDto
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "Room 1", resp["room"].Name)
}

func TestGetRooms_Success(t *testing.T) {
	deps := newTestDeps(t)
	token := "valid-token"

	deps.jwtProvider.On("ParseToken", token).Return(auth.Token{
		UserId: uuid.New(),
		Role:   "user",
	}, nil)
	now := time.Now()
	mockRooms := []domain.Room{
		{Id: uuid.New(), Name: "Room A", CreatedAt: &now},
		{Id: uuid.New(), Name: "Room B", CreatedAt: &now},
	}

	deps.service.On("GetRooms", mock.Anything).Return(mockRooms, nil)

	headers := map[string]string{"Authorization": "Bearer " + token}
	w := makeRequest(deps.router, "GET", "/rooms/list", nil, headers)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]dto.RoomDto
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Len(t, resp["rooms"], 2)
}

func TestGetSlots_MissingDate(t *testing.T) {
	deps := newTestDeps(t)
	token := "valid-token"
	roomId := uuid.New()

	deps.jwtProvider.On("ParseToken", token).Return(auth.Token{
		UserId: uuid.New(),
		Role:   "user",
	}, nil)

	headers := map[string]string{"Authorization": "Bearer " + token}
	url := fmt.Sprintf("/rooms/%s/slots/list", roomId)
	w := makeRequest(deps.router, "GET", url, nil, headers)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing date argument")
}
