package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type BookingRequest struct {
	Id             uuid.UUID
	SlotId         uuid.UUID
	UserId         uuid.UUID
	Status         BookingStatus
	ConferenceLink *string
	CreatedAt      time.Time
}
