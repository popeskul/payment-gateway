package ports

//go:generate mockgen -destination=repository_mock.go -package=ports github.com/popeskul/payment-gateway/internal/core/ports Repositories,MerchantRepository,PaymentRepository,RefundRepository,UserRepository
//go:generate mockgen -destination=service_mock.go -package=ports github.com/popeskul/payment-gateway/internal/core/ports Services,MerchantService,AcquiringBank,PaymentService,RefundService,UserService
//go:generate mockgen -destination=auth_mock.go -package=ports github.com/popeskul/payment-gateway/internal/core/ports AuthConfig,TokenStore,JWTManager,PasswordHasher
//go:generate mockgen -destination=logger_mock.go -package=ports github.com/popeskul/payment-gateway/internal/core/ports Logger
//go:generate mockgen -destination=transaction_mock.go -package=ports github.com/popeskul/payment-gateway/internal/core/ports Transaction
//go:generate mockgen -destination=uuid_generator_mock.go -package=ports github.com/popeskul/payment-gateway/internal/core/ports UUIDGenerator
