package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/popeskul/payment-gateway/internal/core/domain/merchant"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type MerchantRepository struct {
	db            *Database
	uuidGenerator ports.UUIDGenerator
}

func NewMerchantRepository(db *Database, uuidGenerator ports.UUIDGenerator) ports.MerchantRepository {
	return &MerchantRepository{
		db:            db,
		uuidGenerator: uuidGenerator,
	}
}

func (r *MerchantRepository) Create(ctx context.Context, m *merchant.Merchant) error {
	if m.ID == "" {
		m.ID = r.uuidGenerator.Generate()
	}

	query := `
        INSERT INTO merchants (id, name, email, api_key, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.Pool.Exec(ctx, query, m.ID, m.Name, m.Email, m.ApiKey, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create merchant: %v", err)
	}
	return nil
}

func (r *MerchantRepository) GetByID(ctx context.Context, id string) (*merchant.Merchant, error) {
	query := `
		SELECT id, name, email, api_key, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`
	var m merchant.Merchant
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.Email, &m.ApiKey, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("merchant not found")
		}
		return nil, fmt.Errorf("failed to get merchant: %v", err)
	}
	return &m, nil
}

func (r *MerchantRepository) Update(ctx context.Context, m *merchant.Merchant) error {
	query := `
		UPDATE merchants
		SET name = $2, email = $3, api_key = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, m.ID, m.Name, m.Email, m.ApiKey, m.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update merchant: %v", err)
	}
	return nil
}

func (r *MerchantRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM merchants WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete merchant: %v", err)
	}
	return nil
}

func (r *MerchantRepository) List(ctx context.Context, limit, offset int) ([]*merchant.Merchant, error) {
	query := `
		SELECT id, name, email, api_key, created_at, updated_at
		FROM merchants
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list merchants: %v", err)
	}
	defer rows.Close()

	var merchants []*merchant.Merchant
	for rows.Next() {
		var m merchant.Merchant
		err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.ApiKey, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan merchant: %v", err)
		}
		merchants = append(merchants, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating merchants: %v", err)
	}

	return merchants, nil
}

func (r *MerchantRepository) GetByEmail(ctx context.Context, email string) (*merchant.Merchant, error) {
	query := `
		SELECT id, name, email, api_key, created_at, updated_at
		FROM merchants
		WHERE email = $1
	`
	var m merchant.Merchant
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&m.ID, &m.Name, &m.Email, &m.ApiKey, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("merchant not found")
		}
		return nil, fmt.Errorf("failed to get merchant: %v", err)
	}
	return &m, nil
}
