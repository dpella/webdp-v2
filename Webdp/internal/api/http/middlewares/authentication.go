package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"webdp/internal/api/http/repo/postgres"
	"webdp/internal/api/http/services"

	"github.com/golang-jwt/jwt/v4"

	errors "webdp/internal/api/http"
)

type DPContextKey struct {
	Key string
}

const (
	UserContextKey string = "user"
)

func GetTokenAuthentication(tokenRepo postgres.TokenPostgres) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var claims services.JWTTokenClaims
			tokenString, err := ExtracAuthnHeader(r.Header, &claims)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if claims.ExpiresAt < time.Now().Unix() {
				http.Error(w, "your token has expired", http.StatusUnauthorized)
				return
			}

			user := claims.Handle
			savedToken, err := tokenRepo.GetUserToken(user)
			if err != nil {
				http.Error(w, "unauthorized: faulty token", http.StatusUnauthorized)
				return
			}

			if tokenString != savedToken {
				http.Error(w, "unauthorized: faulty token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), DPContextKey{Key: UserContextKey}, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtracAuthnHeader[C jwt.Claims](header http.Header, v C) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("%w: missing header", errors.ErrBadRequest)
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("%w: expected prefix \"Bearer \"", errors.ErrBadRequest)
	}

	auth = strings.Replace(auth, "Bearer ", "", 1)

	_, err := jwt.ParseWithClaims(auth, v, func(t *jwt.Token) (interface{}, error) {
		key := os.Getenv("AUTH_SIGN_KEY")
		if key == "" {
			// under development
			panic("missing key")
		}
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	return auth, v.Valid()

}
