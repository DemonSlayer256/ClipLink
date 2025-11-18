package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	m "ClipLink/models"
)
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Missing or Invalid token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			http.Error(w, "Missing or Invalid token", http.StatusUnauthorized)
			return
		}
		var jwt m.JWT
		if err := json.Unmarshal([]byte(token), &jwt); err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		expirationTime, err := time.Parse(time.RFC3339, jwt.Exp)
		if err != nil || expirationTime.Before(time.Now()) {
			http.Error(w, "Token Expired", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "username", jwt.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}