package postgres_repository

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	"avitoBooking/internal/repository"
	repository_errors "avitoBooking/internal/repository/erorrs"
	"avitoBooking/internal/repository/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type authRepository struct {
	TxManager repository.TxManager
}

func NewAuthRepository(
	txManager repository.TxManager,
) repository.AuthRepository {
	return &authRepository{TxManager: txManager}
}

func (r *authRepository) Register(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()
	sqlQuery := `
	INSERT INTO users(role,email,password_hash)
	VALUES($1,$2,$3)
	RETURNING id,role,email,password_hash,created_at;
	`
	var model models.UserModel
	err := r.TxManager.GetExecutor(ctx).QueryRow(
		ctx,
		sqlQuery,
		user.Role,
		user.Email,
		user.Password,
	).Scan(
		&model.Id,
		&model.Role,
		&model.Email,
		&model.Password,
		&model.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.User{}, core_errors.ErrEmailExists
			}
		}
		return domain.User{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return models.UserModelToDomain(model), nil
}
func (r *authRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TxManager.GetTimeout())
	defer cancel()

	sqlQuery := `
	SELECT id,role,email,password_hash,created_at
	FROM users
	WHERE email = $1;
	`

	var user models.UserModel
	err := r.TxManager.GetExecutor(ctx).QueryRow(ctx, sqlQuery, email).Scan(
		&user.Id,
		&user.Role,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("failed to get user by email: %w", repository_errors.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return domain.User{
		Id:        user.Id,
		Role:      user.Role,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}, nil
}
