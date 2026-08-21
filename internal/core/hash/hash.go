package core_hash

import (
	core_config "avitoBooking/internal/core/config"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	Cost int
}
type HasherInterface interface {
	HashPassword(password []byte) ([]byte, error)
	Compare(password []byte, passwordHash []byte) error
}

func NewHasher(cfg core_config.Config) *Hasher {
	return &Hasher{
		Cost: cfg.BcryptCost,
	}
}

func (h *Hasher) HashPassword(password []byte) ([]byte, error) {

	hash, err := bcrypt.GenerateFromPassword(password, h.Cost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash to password: %w", err)
	}
	return hash, nil
}
func (h *Hasher) Compare(password []byte, passwordHash []byte) error {
	return bcrypt.CompareHashAndPassword(passwordHash, password)
}
