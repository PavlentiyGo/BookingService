package domain

import (
	core_errors "avitoBooking/internal/core/errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Token string

type User struct {
	Id        uuid.UUID
	Role      string
	Email     string
	Password  []byte
	CreatedAt time.Time
}

func (u *User) Validate() error {
	if u.Email != "" {
		emailReg := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailReg.MatchString(u.Email) {
			return fmt.Errorf("invalid email: %w",
				core_errors.ErrInvalidRequest,
			)
		}
	}
	if len(u.Password) < 3 {
		return fmt.Errorf("too small password: %w",
			core_errors.ErrInvalidRequest,
		)
	}

	return nil
}
