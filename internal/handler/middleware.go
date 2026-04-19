package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Snuslyk/schedule-backend-go/internal/config"
	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "authUser"

type AuthUser struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			id, _ := claims["id"].(string)
			var roles []string
			if rawRoles, ok := claims["roles"].([]any); ok {
				for _, r := range rawRoles {
					if str, ok := r.(string); ok {
						roles = append(roles, str)
					}
				}
			}

			user := &AuthUser{
				ID:    id,
				Roles: roles,
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthUser(r *http.Request) *AuthUser {
	if u, ok := r.Context().Value(userContextKey).(*AuthUser); ok {
		return u
	}
	return nil
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := GetAuthUser(r)
		if u == nil {
			Unauthorized(w, "Unauthorized")
			return
		}
		next(w, r)
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := GetAuthUser(r)
		if u == nil {
			Unauthorized(w, "Unauthorized")
			return
		}

		isAdmin := false
		for _, role := range u.Roles {
			if role == string(domain.RoleAdmin) {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			Forbidden(w, "Forbidden resource")
			return
		}

		next(w, r)
	}
}
