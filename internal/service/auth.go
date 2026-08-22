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

func (s *Service) Login(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, err
	}

	userFromDb, err := s.authRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	if err = s.hasher.Compare(user.Password, userFromDb.Password); err != nil {
		return domain.User{}, err
	}
	return userFromDb, nil
}
