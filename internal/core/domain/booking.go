package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	Id             uuid.UUID
	SlotId         uuid.UUID
	UserId         uuid.UUID
	Status         string
	ConferenceLink *string
	CreatedAt      time.Time
}
