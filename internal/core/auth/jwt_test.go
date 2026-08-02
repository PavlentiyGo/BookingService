package auth

import (
	core_errors "avitoBooking/internal/core/errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwt_Success(t *testing.T) {
	jwtConfig := JwtTokenConfig{
		lifeTime:   5 * time.Second,
		signingKey: "secret_key",
	}
	userId, err := uuid.NewUUID()
	require.NoError(t, err)

	token, err := jwtConfig.NewToken("user", userId)
	require.NoError(t, err)

	parsedToken, err := jwtConfig.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, parsedToken.UserId, userId)
	assert.Equal(t, parsedToken.Role, "user")
}

func TestJwt_NegativeLifeTime(t *testing.T) {
	jwtConfig := JwtTokenConfig{
		lifeTime:   -5 * time.Second,
		signingKey: "secret_key",
	}
	userId, err := uuid.NewUUID()
	require.NoError(t, err)

	token, err := jwtConfig.NewToken("user", userId)
	require.NoError(t, err)

	_, err = jwtConfig.ParseToken(token)
	assert.ErrorIs(t, err, core_errors.ErrExpiredToken)
}
func TestJwt_InvalidToken(t *testing.T) {
	jwtConfig := JwtTokenConfig{
		lifeTime:   5 * time.Second,
		signingKey: "secret_key",
	}
	_, err := jwtConfig.ParseToken("jwt_token")
	assert.ErrorIs(t, err, core_errors.ErrInvalidRequest)
}
func TestJwt_DifferentKey(t *testing.T) {
	jwtConfig := JwtTokenConfig{
		lifeTime:   5 * time.Second,
		signingKey: "secret_key",
	}
	userId, err := uuid.NewUUID()
	require.NoError(t, err)

	token, err := jwtConfig.NewToken("user", userId)
	require.NoError(t, err)
	jwtConfig.signingKey = "different_key"
	_, err = jwtConfig.ParseToken(token)
	assert.Error(t, err)
}
