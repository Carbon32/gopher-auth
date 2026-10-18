package auth

import (
	"carb-backend/service/types"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const UserIdKey = "id"
const SessionIDKey = "session_id"

func RequireAuth(secret []byte, check types.SessionChecker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims := new(types.AccessClaims)
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || token == nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil || userID <= 0 || claims.SessionID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		active, err := check(r.Context(), claims.SessionID)
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		if !active {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, userID)
		ctx = context.WithValue(ctx, SessionIDKey, claims.SessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
