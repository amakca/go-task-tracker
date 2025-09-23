package v1

import (
	service "go-task-tracker/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(r chi.Router, services *service.Services) {
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(SlogRequestContext)
	r.Use(SlogAccessLogger)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	// r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", func(cr chi.Router) {
		newAuthRoutes(cr, services.Auth)
	})

	authMiddleware := &AuthMiddleware{services.Auth}
	r.Route("/api/v1", func(api chi.Router) {
		api.Use(authMiddleware.UserIdentity)
		// Остальные группы хендлеров
	})
}
