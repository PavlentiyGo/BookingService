package auth

import (
	"avitoBooking/internal/core/errors"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	jwt.MapClaims
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type JwtProvider interface {
	NewToken(role string, userId uuid.UUID) (string, error)
	ParseToken(tokenString string) (Token, error)
}

type JwtTokenConfig struct {
	lifeTime   time.Duration
	signingKey string
}

func NewJwtConfig(
	lifeTime time.Duration,
	key string,
) *JwtTokenConfig {
	return &JwtTokenConfig{
		lifeTime:   lifeTime,
		signingKey: key,
	}
}

func (j *JwtTokenConfig) NewToken(role string, userId uuid.UUID) (string, error) {
	claims := Token{
		MapClaims: jwt.MapClaims{
			"exp": jwt.NewNumericDate(time.Now().Add(j.lifeTime)),
		},
		UserId: userId,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	newToken, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt token: %w: %w", core_errors.ErrInternalError, err)
	}
	return newToken, nil
}

func (j *JwtTokenConfig) ParseToken(tokenString string) (Token, error) {

	tokenClaims := Token{}
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("wrong signing method: %v: %w", token.Header["alg"], core_errors.ErrInvalidRequest)
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return Token{}, fmt.Errorf("invalid token: %w", core_errors.ErrExpiredToken)
		}
		return Token{}, fmt.Errorf("parse token: %w: %w", err, core_errors.ErrInvalidRequest)
	}
	if !token.Valid {
		return Token{}, fmt.Errorf("invalid token: %w", core_errors.ErrInvalidRequest)
	}
	return tokenClaims, nil
}
