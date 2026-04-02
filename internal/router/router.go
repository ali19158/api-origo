package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/online-shop/internal/handler"
	"github.com/online-shop/internal/middleware"
)

func New(
	jwtSecret string,
	userH *handler.UserHandler,
	productH *handler.ProductHandler,
	//orderH *handler.OrderHandler,
	categoryH *handler.CategoryHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth
		r.Post("/auth/register", userH.Register)
		r.Post("/auth/login", userH.Login)

		// Public product & category listing
		r.Get("/products", productH.List)
		r.Get("/products/{id}", productH.GetByID)
		r.Get("/categories", categoryH.List)
		r.Get("/categories/{id}", categoryH.GetByID)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))

			// Orders (customer)
			//r.Post("/orders", orderH.Create)
			//r.Get("/orders/my", orderH.ListMy)
			//r.Get("/orders/{id}", orderH.GetByID)

			// Admin routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdmin)

				r.Post("/products", productH.Create)
				r.Put("/products/{id}", productH.Update)
				r.Delete("/products/{id}", productH.Delete)

				r.Post("/categories", categoryH.Create)
				r.Put("/categories/{id}", categoryH.Update)
				r.Delete("/categories/{id}", categoryH.Delete)

				//r.Put("/orders/{id}/status", orderH.UpdateStatus)
			})
		})
	})

	return r
}
