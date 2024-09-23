package api

import (
	"github.com/popeskul/payment-gateway/internal/infrastructure/metrics"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/popeskul/payment-gateway/internal/api/handlers"
	customMiddleware "github.com/popeskul/payment-gateway/internal/api/middleware"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type Router struct {
	router  *chi.Mux
	handler *handlers.Handler
}

func NewRouter(services ports.Services, logger ports.Logger, jwtManager ports.JWTManager) *Router {
	r := &Router{
		router:  chi.NewRouter(),
		handler: handlers.NewHandler(services, logger, jwtManager),
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	r.router.Use(middleware.RequestID)
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Recoverer)
	r.router.Use(customMiddleware.Cors)
	r.router.Use(customMiddleware.MetricsMiddleware)

	r.router.Handle("/metrics", metrics.MetricsHandler())

	r.router.Route("/api/v1", func(router chi.Router) {
		// Public authentication routes
		router.Post("/register", r.handler.Register)
		router.Post("/login", r.handler.Login)
		router.Post("/refresh", r.handler.RefreshToken)

		// Protected routes
		router.Group(func(router chi.Router) {
			router.Use(customMiddleware.Auth(r.handler.JWTManager))

			// User routes
			router.Post("/logout", r.handler.Logout)
			router.Get("/profile", r.handler.GetProfile)
			router.Put("/profile", r.handler.UpdateProfile)
			router.Post("/change-password", r.handler.ChangePassword)

			// Merchant routes
			router.Post("/merchants", r.handler.CreateMerchant)
			router.Get("/merchants/{id}", r.handler.GetMerchant)
			router.Put("/merchants/{id}", r.handler.UpdateMerchant)
			router.Delete("/merchants/{id}", r.handler.DeleteMerchant)
			router.Get("/merchants", r.handler.ListMerchants)

			// Payment routes
			router.Post("/payments", r.handler.CreatePayment)
			router.Get("/payments/{id}", r.handler.GetPayment)
			router.Get("/payments", r.handler.ListPayments)
			router.Post("/payments/{id}/process", r.handler.ProcessPayment)

			// Refund routes
			router.Post("/refunds", r.handler.CreateRefund)
			router.Get("/refunds/{id}", r.handler.GetRefund)
			router.Get("/refunds", r.handler.ListRefunds)
			router.Post("/refunds/{id}/process", r.handler.ProcessRefund)
		})
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
