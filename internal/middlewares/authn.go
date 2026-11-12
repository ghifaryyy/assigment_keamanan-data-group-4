package middlewares

import (
	"aaa_bekend/internal/lmodel"
	"aaa_bekend/internal/auth"
	"net/http"
	"strings"
)

func AuthN(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// TODO
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		
		if !strings.HasPrefix(authHeader, "Bearer"){
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// TODO
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if claims.UserID <= 0 || claims.Username == "" || claims.Role == "" {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userCtx := lmodel.UserContext{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx := lmodel.ContextWithUser(r.Context(), userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
