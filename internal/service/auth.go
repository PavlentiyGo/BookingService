package service

import (
	"avitoBooking/internal/core/domain"
	"context"
	"fmt"
)

func (s *Service) Register(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	passwordHash, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = passwordHash
	newUser, err := s.authRepo.Register(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to register user: %w", err)
	}
	return newUser, nil
}
