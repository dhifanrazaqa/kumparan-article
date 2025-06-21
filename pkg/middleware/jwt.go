package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/dhifanrazaqa/kumparan-article/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsContextKey contextKey = "userClaims"

func JWT(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "Authorization format must be 'Bearer <token>'")
			return
		}
		tokenString := parts[1]

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, "Token is invalid or expired")
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}