package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type PostgresTokenStore struct {
	pool *pgxpool.Pool
}

func NewPostgresTokenStore(pool *pgxpool.Pool) ports.TokenStore {
	return &PostgresTokenStore{pool: pool}
}

func (s *PostgresTokenStore) StoreRefreshToken(userID string, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx,
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, time.Now().Add(7*24*time.Hour)) // token valid for 7 days
	return err
}

func (s *PostgresTokenStore) DeleteRefreshToken(userID string, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx,
		"DELETE FROM refresh_tokens WHERE user_id = $1 AND token = $2",
		userID, token)
	return err
}

func (s *PostgresTokenStore) IsRefreshTokenValid(userID string, token string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := s.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1 AND token = $2 AND expires_at > NOW()",
		userID, token).Scan(&count)

	return err == nil && count > 0
}
