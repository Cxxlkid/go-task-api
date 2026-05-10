package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Récupérer le header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Format attendu : "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			// Valider le token
			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				respondError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// Extraire le user_id et l'injecter dans le contexte
			claims := token.Claims.(jwt.MapClaims)
			userID := claims["user_id"].(string)
			ctx := context.WithValue(r.Context(), "user_id", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}