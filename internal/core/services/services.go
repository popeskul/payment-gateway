package services

import (
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type Services struct {
	merchantService ports.MerchantService
	paymentService  ports.PaymentService
	refundService   ports.RefundService
	userService     ports.UserService
}

func NewServices(merchantService ports.MerchantService, paymentService ports.PaymentService, refundService ports.RefundService, userService ports.UserService) *Services {
	return &Services{
		merchantService: merchantService,
		paymentService:  paymentService,
		refundService:   refundService,
		userService:     userService,
	}
}

func (s *Services) Merchants() ports.MerchantService {
	return s.merchantService
}

func (s *Services) Payments() ports.PaymentService {
	return s.paymentService
}

func (s *Services) Refunds() ports.RefundService {
	return s.refundService
}

func (s *Services) Users() ports.UserService {
	return s.userService
}
