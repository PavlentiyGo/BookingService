package core_middleware

import (
	"context"

	"github.com/google/uuid"
)

type ctxUserIdKey = struct{}

type ctxUserRoleKey = struct{}

func ContextWithUserId(
	ctx context.Context,
	userId uuid.UUID,
) context.Context {
	return context.WithValue(ctx, ctxUserIdKey{}, userId)
}
func ContextWithUserRole(
	ctx context.Context,
	role string,
) context.Context {
	return context.WithValue(ctx, ctxUserRoleKey{}, role)
}
func UserIdFromContext(
	ctx context.Context,
) uuid.UUID {
	return ctx.Value(ctxUserIdKey{}).(uuid.UUID)
}
func UserRoleFromContext(
	ctx context.Context,
) string {
	return ctx.Value(ctxUserRoleKey{}).(string)
}
