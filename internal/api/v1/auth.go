package v1

import (
	"github.com/go-chi/chi/v5"
)

func registerAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		// r.Post("/signup", handleSignup)
		// r.Post("/login", handleLogin)
	})
}

type signupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
}

// func (a *API) handleSignup(w http.ResponseWriter, r *http.Request) {
// 	var req signupRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}
// 	if req.Username == "" || req.Email == "" || len(req.Password) < 6 {
// 		http.Error(w, "invalid fields", http.StatusBadRequest)
// 		return
// 	}
// 	hash, err := auth.HashPassword(req.Password)
// 	if err != nil {
// 		log.Error().Err(err).Msg("hash error")
// 		http.Error(w, "server error", http.StatusInternalServerError)
// 		return
// 	}
// 	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
// 	defer cancel()
// 	row := a.DB.Pool.QueryRow(ctx, `INSERT INTO app_user (username,email,password_hash) VALUES ($1,$2,$3) RETURNING id`, req.Username, strings.ToLower(req.Email), hash)
// 	var userID string
// 	if err := row.Scan(&userID); err != nil {
// 		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
// 			http.Error(w, "user exists", http.StatusConflict)
// 			return
// 		}
// 		log.Error().Err(err).Msg("insert user error")
// 		http.Error(w, "server error", http.StatusInternalServerError)
// 		return
// 	}
// 	tok, err := auth.CreateAccessToken(userID)
// 	if err != nil {
// 		log.Error().Err(err).Msg("jwt error")
// 		http.Error(w, "server error", http.StatusInternalServerError)
// 		return
// 	}
// 	writeJSON(w, http.StatusOK, authResponse{AccessToken: tok})
// }

// func (a *API) handleLogin(w http.ResponseWriter, r *http.Request) {
// 	var req loginRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}
// 	if req.Username == "" || req.Password == "" {
// 		http.Error(w, "invalid fields", http.StatusBadRequest)
// 		return
// 	}
// 	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
// 	defer cancel()
// 	var userID, passwordHash string
// 	err := a.DB.Pool.QueryRow(ctx, `SELECT id, password_hash FROM app_user WHERE username=$1`, req.Username).Scan(&userID, &passwordHash)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			http.Error(w, "invalid credentials", http.StatusUnauthorized)
// 			return
// 		}
// 		log.Error().Err(err).Msg("query user error")
// 		http.Error(w, "server error", http.StatusInternalServerError)
// 		return
// 	}
// 	ok, err := auth.VerifyPassword(req.Password, passwordHash)
// 	if err != nil || !ok {
// 		http.Error(w, "invalid credentials", http.StatusUnauthorized)
// 		return
// 	}
// 	tok, err := auth.CreateAccessToken(userID)
// 	if err != nil {
// 		log.Error().Err(err).Msg("jwt error")
// 		http.Error(w, "server error", http.StatusInternalServerError)
// 		return
// 	}
// 	writeJSON(w, http.StatusOK, authResponse{AccessToken: tok})
// }

// func handleHealth(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	_, _ = w.Write([]byte("ok"))
// }

// func writeJSON(w http.ResponseWriter, code int, v any) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	_ = json.NewEncoder(w).Encode(v)
// }
