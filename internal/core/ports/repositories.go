package ports

import (
	"context"

	"github.com/popeskul/payment-gateway/internal/core/domain/merchant"
	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/domain/refund"
	"github.com/popeskul/payment-gateway/internal/core/domain/user"
)

type Repositories interface {
	Merchants() MerchantRepository
	Payments() PaymentRepository
	Refunds() RefundRepository
	Users() UserRepository
}

type MerchantRepository interface {
	Create(ctx context.Context, m *merchant.Merchant) error
	GetByID(ctx context.Context, id string) (*merchant.Merchant, error)
	Update(ctx context.Context, m *merchant.Merchant) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*merchant.Merchant, error)
	GetByEmail(ctx context.Context, email string) (*merchant.Merchant, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, p *payment.Payment) error
	GetByID(ctx context.Context, id string) (*payment.Payment, error)
	Update(ctx context.Context, p *payment.Payment) error
	List(ctx context.Context, merchantID string, limit, offset int) ([]*payment.Payment, error)
	UpdateStatus(ctx context.Context, id string, status payment.PaymentStatus) error
}

type RefundRepository interface {
	Create(ctx context.Context, r *refund.Refund) error
	GetByID(ctx context.Context, id string) (*refund.Refund, error)
	Update(ctx context.Context, r *refund.Refund) error
	List(ctx context.Context, paymentID string, limit, offset int) ([]*refund.Refund, error)
	UpdateStatus(ctx context.Context, id string, status refund.RefundStatus) error
	UpdateWithTransaction(ctx context.Context, r *refund.Refund, p *payment.Payment) error
}

type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
}
