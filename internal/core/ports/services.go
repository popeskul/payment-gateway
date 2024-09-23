package ports

import (
	"context"
	"time"

	"github.com/popeskul/payment-gateway/internal/core/domain/merchant"
	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/domain/refund"
	"github.com/popeskul/payment-gateway/internal/core/domain/user"
)

type Services interface {
	Merchants() MerchantService
	Payments() PaymentService
	Refunds() RefundService
	Users() UserService
}

type MerchantService interface {
	CreateMerchant(ctx context.Context, m *merchant.Merchant) error
	GetMerchant(ctx context.Context, id string) (*merchant.Merchant, error)
	UpdateMerchant(ctx context.Context, m *merchant.Merchant) error
	DeleteMerchant(ctx context.Context, id string) error
	ListMerchants(ctx context.Context, limit, offset int) ([]*merchant.Merchant, error)
}

type AcquiringBank interface {
	ProcessPayment(ctx context.Context, p *payment.Payment) error
	ProcessRefund(ctx context.Context, r *refund.Refund) error
	SetProcessingDelay(delay time.Duration)
	SetFailureRate(rate float64)
}

type PaymentService interface {
	CreatePayment(ctx context.Context, p *payment.Payment) error
	GetPayment(ctx context.Context, id string) (*payment.Payment, error)
	UpdatePayment(ctx context.Context, p *payment.Payment) error
	ListPayments(ctx context.Context, merchantID string, limit, offset int) ([]*payment.Payment, error)
	ProcessPayment(ctx context.Context, paymentID string) error
}

type RefundService interface {
	CreateRefund(ctx context.Context, r *refund.Refund) error
	GetRefund(ctx context.Context, id string) (*refund.Refund, error)
	UpdateRefund(ctx context.Context, r *refund.Refund) error
	ListRefunds(ctx context.Context, paymentID string, limit, offset int) ([]*refund.Refund, error)
	ProcessRefund(ctx context.Context, refundID string) error
}

type UserService interface {
	Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error)
	Login(ctx context.Context, req *user.LoginRequest) (*user.User, error)
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	UpdateProfile(ctx context.Context, id string, req *user.UpdateProfileRequest) (*user.User, error)
	ChangePassword(ctx context.Context, id string, req *user.ChangePasswordRequest) error
}
