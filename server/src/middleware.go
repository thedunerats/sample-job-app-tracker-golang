package main

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is a custom type for context keys
type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware validates JWT tokens and adds user info to request context
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves user claims from request context
func GetUserFromContext(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(userContextKey).(*Claims)
	return claims, ok
}
