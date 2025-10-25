package rest

import (
	"hyperlocal/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// SetupRouter configures all routes for the application
func SetupRouter(handler *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler())

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes - no middleware required
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handler.V1.Register)
			r.Post("/login", handler.V1.Login)
			r.Post("/refresh", handler.V1.RefreshToken)
		})

		// Protected routes - require authentication
		r.Group(func(r chi.Router) {
			r = r.With(AuthMiddlewareFunc(handler.V1.Service))

			// Posts
			r.Route("/posts", func(r chi.Router) {
				r.With(RateLimiterMiddleware).Post("/", handler.V1.CreatePost)
				r.Get("/", handler.V1.GetNearbyPosts)

				// Post interactions
				r.Post("/{id}/upvote", handler.V1.UpvotePost)
				r.Post("/{id}/downvote", handler.V1.DownvotePost)
				r.Post("/{id}/report", handler.V1.ReportPost)

				// Comments
				r.Post("/{id}/comments", handler.V1.CreateComment)
				r.Get("/{id}/comments", handler.V1.GetComments)
			})

			// Admin routes - require admin role
			r.Route("/admin", func(r chi.Router) {
				r.Use(AdminMiddleware)

				r.Get("/flagged", handler.V1.GetFlaggedPosts)
				r.Delete("/posts/{id}", handler.V1.DeletePost)
				r.Patch("/users/{id}/ban", handler.V1.BanUser)
			})
		})
	})

	return r
}
