package rest

import (
	"hyperlocal/internal/handlers"

	"github.com/gorilla/mux"
)

// NewRouter returns a new router instance with configured routes
func NewRouter(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	return router
}
