package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/domain/refund"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type RefundRepository struct {
	db            *Database
	uuidGenerator ports.UUIDGenerator
}

func NewRefundRepository(db *Database, uuidGenerator ports.UUIDGenerator) ports.RefundRepository {
	return &RefundRepository{
		db:            db,
		uuidGenerator: uuidGenerator,
	}
}

func (r *RefundRepository) Create(ctx context.Context, ref *refund.Refund) error {
	if ref.ID == "" {
		ref.ID = r.uuidGenerator.Generate()
	}

	query := `
		INSERT INTO refunds (id, payment_id, amount, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		ref.ID, ref.PaymentID, ref.Amount, ref.Reason, ref.Status, ref.CreatedAt, ref.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create refund: %v", err)
	}
	return nil
}

func (r *RefundRepository) GetByID(ctx context.Context, id string) (*refund.Refund, error) {
	query := `
		SELECT id, payment_id, amount, reason, status, created_at, updated_at
		FROM refunds
		WHERE id = $1
	`
	var ref refund.Refund
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&ref.ID, &ref.PaymentID, &ref.Amount, &ref.Reason, &ref.Status, &ref.CreatedAt, &ref.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("refund not found")
		}
		return nil, fmt.Errorf("failed to get refund: %v", err)
	}
	return &ref, nil
}

func (r *RefundRepository) Update(ctx context.Context, ref *refund.Refund) error {
	query := `
		UPDATE refunds
		SET payment_id = $2, amount = $3, reason = $4, status = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query,
		ref.ID, ref.PaymentID, ref.Amount, ref.Reason, ref.Status, ref.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update refund: %v", err)
	}
	return nil
}

func (r *RefundRepository) List(ctx context.Context, paymentID string, limit, offset int) ([]*refund.Refund, error) {
	query := `
		SELECT id, payment_id, amount, reason, status, created_at, updated_at
		FROM refunds
		WHERE payment_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, paymentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list refunds: %v", err)
	}
	defer rows.Close()

	var refunds []*refund.Refund
	for rows.Next() {
		var ref refund.Refund
		err := rows.Scan(&ref.ID, &ref.PaymentID, &ref.Amount, &ref.Reason, &ref.Status, &ref.CreatedAt, &ref.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan refund: %v", err)
		}
		refunds = append(refunds, &ref)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating refunds: %v", err)
	}

	return refunds, nil
}

func (r *RefundRepository) UpdateStatus(ctx context.Context, id string, status refund.RefundStatus) error {
	query := `
		UPDATE refunds
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update refund status: %v", err)
	}
	return nil
}

func (r *RefundRepository) UpdateWithTransaction(ctx context.Context, ref *refund.Refund, p *payment.Payment) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	refundQuery := `
		UPDATE refunds
		SET status = $2, updated_at = $3
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, refundQuery, ref.ID, ref.Status, ref.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update refund in transaction: %v", err)
	}

	paymentQuery := `
		UPDATE payments
		SET amount = $2, updated_at = $3
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, paymentQuery, p.ID, p.Amount, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update payment in transaction: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
