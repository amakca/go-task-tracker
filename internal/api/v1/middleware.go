package v1

import (
	"context"
	"go-task-tracker/internal/services/contracts"
	"go-task-tracker/pkg/logctx"
	"net/http"
	"strings"
)

const (
	userIdCtx = "userId"
)

type AuthMiddleware struct {
	authService contracts.Auth
}

func (h *AuthMiddleware) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logctx.FromContext(r.Context())

		token, ok := bearerToken(r)
		if !ok {
			log.Warn("AuthMiddleware.UserIdentity: bearerToken", "error", ErrInvalidAuthHeader)
			newErrorResponseHTTP(w, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return
		}

		userId, err := h.authService.ParseToken(token)
		if err != nil {
			log.Warn("AuthMiddleware.UserIdentity: ParseToken", "err", err)
			newErrorResponseHTTP(w, http.StatusUnauthorized, ErrCannotParseToken.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userIdCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
