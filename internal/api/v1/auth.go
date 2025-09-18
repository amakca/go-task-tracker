package v1

import (
	"context"
	"encoding/json"
	"go-task-tracker/internal/services/contracts"
	"go-task-tracker/internal/services/users"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type authRoutes struct {
	authService contracts.Auth
}

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerAuthRoutes(r chi.Router, authService contracts.Auth) {
	a := &authRoutes{
		authService: authService,
	}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", a.handleSignup)
		r.Post("/login", a.handleLogin)
	})
}

func (a *authRoutes) handleSignup(w http.ResponseWriter, r *http.Request) {

	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Username == "" || len(req.Password) < 6 {
		http.Error(w, "invalid fields", http.StatusBadRequest)
		return
	}

	ctx := context.TODO()
	id, err := a.authService.CreateUser(ctx, contracts.AuthCreateUserInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Error().Err(err)
		if err == users.ErrUserAlreadyExists {
			http.Error(w, "server error", http.StatusBadRequest)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	writeJSON(w, http.StatusCreated, response{Id: id})
}

func (a *authRoutes) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "invalid fields", http.StatusBadRequest)
		return
	}

	ctx := context.TODO()
	token, err := a.authService.GenerateToken(ctx, contracts.AuthGenerateTokenInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Error().Err(err)
		if err == users.ErrUserNotFound {
			http.Error(w, "server error", http.StatusBadRequest)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	writeJSON(w, http.StatusCreated, response{Token: token})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
