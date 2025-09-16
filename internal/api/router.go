package api

import (
	"go-task-tracker/internal/repo"

	"github.com/go-chi/chi/v5"
)

type API struct {
	DB *repo.DB
}

func NewAPI(db *repo.DB) *API {
	return &API{DB: db}
}

func (a *API) RegisterRoutes(r chi.Router) {
	r.Get("/healthz", handleHealth)
	a.registerAuthRoutes(r)
}
