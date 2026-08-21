package service

import (
	core_errors "avitoBooking/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ConferenceService struct {
	available bool
}

func NewConferenceService(
	isAvailable bool,
) *ConferenceService {
	return &ConferenceService{
		available: isAvailable,
	}
}

func (s *ConferenceService) CreateConferenceLink(
	_ context.Context,
	slotId uuid.UUID,
) (string, error) {
	if !s.available {
		return "", core_errors.ErrConferenceServiceUnavailable
	}
	link := fmt.Sprintf("https://conference/%s", slotId.String())

	return link, nil
}
func (s *ConferenceService) CancelConference(_ context.Context, _ string) error {
	if !s.available {
		return core_errors.ErrConferenceServiceUnavailable
	}

	return nil
}
