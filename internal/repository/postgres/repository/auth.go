package postgres_repository

import (
	"avitoBooking/internal/core/domain"
	core_errors "avitoBooking/internal/core/errors"
	"avitoBooking/internal/repository"
	"avitoBooking/internal/repository/models"
	"context"
	"errors"
	"fmt"

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
