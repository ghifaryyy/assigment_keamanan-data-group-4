package middlewares

import (
	"aaa_bekend/internal/auth"
	"aaa_bekend/internal/lmodel"
	"net/http"
	"strings"
)

func AuthN(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// TODO
		if claims.UserID <= 0 || claims.Username == "" || claims.Role == "" {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		// TODO
		userCtx := lmodel.UserContext{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx := lmodel.ContextWithUser(r.Context(), userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
