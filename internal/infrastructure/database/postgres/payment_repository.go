package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type PaymentRepository struct {
	db            *Database
	uuidGenerator ports.UUIDGenerator
}

func NewPaymentRepository(db *Database, uuidGenerator ports.UUIDGenerator) ports.PaymentRepository {
	return &PaymentRepository{
		db:            db,
		uuidGenerator: uuidGenerator,
	}
}

func (r *PaymentRepository) Create(ctx context.Context, p *payment.Payment) error {
	if p.ID == "" {
		p.ID = r.uuidGenerator.Generate()
	}

	query := `
		INSERT INTO payments (id, merchant_id, amount, currency, status, payment_method, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Pool.Exec(ctx, query, p.ID, p.MerchantID, p.Amount, p.Currency, p.Status, p.PaymentMethod, p.Description, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create payment: %v", err)
	}
	return nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*payment.Payment, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, payment_method, description, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	var p payment.Payment
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.MerchantID, &p.Amount, &p.Currency, &p.Status, &p.PaymentMethod, &p.Description, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}
	return &p, nil
}

func (r *PaymentRepository) Update(ctx context.Context, p *payment.Payment) error {
	query := `
		UPDATE payments
		SET merchant_id = $2, amount = $3, currency = $4, status = $5, payment_method = $6, description = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query,
		p.ID, p.MerchantID, p.Amount, p.Currency, p.Status, p.PaymentMethod, p.Description, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update payment: %v", err)
	}
	return nil
}

func (r *PaymentRepository) List(ctx context.Context, merchantID string, limit, offset int) ([]*payment.Payment, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, payment_method, description, created_at, updated_at
		FROM payments
		WHERE merchant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, merchantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %v", err)
	}
	defer rows.Close()

	var payments []*payment.Payment
	for rows.Next() {
		var p payment.Payment
		err := rows.Scan(&p.ID, &p.MerchantID, &p.Amount, &p.Currency, &p.Status, &p.PaymentMethod, &p.Description, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %v", err)
		}
		payments = append(payments, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %v", err)
	}

	return payments, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id string, status payment.PaymentStatus) error {
	query := `
		UPDATE payments
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}
	return nil
}
