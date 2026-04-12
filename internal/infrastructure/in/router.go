package in

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimdware "github.com/go-chi/chi/v5/middleware"

	"wishlist/internal/infrastructure/in/handlers"
	"wishlist/internal/infrastructure/in/middleware"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	wishlistHandler *handlers.WishlistHandler,
	itemHandler *handlers.ItemHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimdware.Logger)
	r.Use(chimdware.Recoverer)
	r.Use(chimdware.Timeout(15 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Route("/wishlists", func(r chi.Router) {
			r.Use(authMiddleware.Handle)

			r.Post("/", wishlistHandler.Create)
			r.Get("/", wishlistHandler.GetAll)
			r.Get("/{id}", wishlistHandler.GetByID)
			r.Put("/{id}", wishlistHandler.Update)
			r.Delete("/{id}", wishlistHandler.Delete)
		})

		r.Route("/items", func(r chi.Router) {
			r.Use(authMiddleware.Handle)

			r.Post("/{wishlist_id}", itemHandler.Create)
			r.Get("/{wishlist_id}", itemHandler.GetAllByWishlist)
			r.Put("/{id}", itemHandler.Update)
			r.Delete("/{id}", itemHandler.Delete)
		})

		r.Get("/public/{token}", wishlistHandler.GetPublic)
		r.Post("/public/{token}/items/{id}/reserve", itemHandler.Reserve)
	})

	return r
}
