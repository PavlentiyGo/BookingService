package postgres

import (
	core_config "avitoBooking/internal/core/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg core_config.DatabaseConfig) (*pgxpool.Pool, error) {

	connString := getConnectionString(cfg)

	newPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}
	if err = newPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping created pool: %w", err)
	}
	return newPool, nil
}

func PoolShutdown(ctx context.Context, pool *pgxpool.Pool) error {

	closed := make(chan bool)

	go func() {
		pool.Close()
		closed <- true
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to close pool before main context done")
	case <-closed:
		return nil
	}
}

func getConnectionString(cfg core_config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=UTC",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}
