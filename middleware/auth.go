package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kizoukun/codingtest/helper"
	"github.com/kizoukun/codingtest/repository"
	"github.com/kizoukun/codingtest/web"
)

func JWTMiddleware(next http.Handler) http.Handler {
	secret := []byte(os.Getenv("JWT_PRIVATE_KEY"))

	var response web.ResponseHttp

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			response.StatusCode = http.StatusUnauthorized
			response.Message = "missing auth header"

			w.WriteHeader(response.StatusCode)
			json.NewEncoder(w).Encode(response)
			return
		}
		// expecting "Bearer <token>"
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			response.StatusCode = http.StatusUnauthorized
			response.Message = "invalid auth header"

			w.WriteHeader(response.StatusCode)
			json.NewEncoder(w).Encode(response)
			return
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			response.StatusCode = http.StatusUnauthorized
			response.Message = "invalid token"

			w.WriteHeader(response.StatusCode)
			json.NewEncoder(w).Encode(response)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.StatusCode = http.StatusUnauthorized
			response.Message = "invalid token claims"

			w.WriteHeader(response.StatusCode)
			json.NewEncoder(w).Encode(response)
			return
		}

		userRepo := repository.NewUserRepository()

		email, err := claims.GetSubject()
		if err != nil {
			response.StatusCode = http.StatusUnauthorized
			response.Message = "invalid token subject"

			w.WriteHeader(response.StatusCode)
			json.NewEncoder(w).Encode(response)
			return
		}

		user, err := userRepo.GetUserByEmail(email)

		if err != nil || user == nil {
			response.StatusCode = http.StatusUnauthorized
			response.Message = "user not found"

			w.WriteHeader(response.StatusCode)
			json.NewEncoder(w).Encode(response)
			return
		}

		// attach user to request context
		ctx := context.WithValue(r.Context(), helper.UserCtxKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
