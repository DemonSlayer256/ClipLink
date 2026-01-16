package middleware

import (
    "context"
    "ClipLink/configs"
    "net/http"
    "strings"
    "time"
    jwt "github.com/golang-jwt/jwt/v4"
    m "ClipLink/models"
)

var signingKey = []byte(configs.LoadEnv("SECURE_KEY")) // You should load this from an environment variable

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if auth == "" {
            http.Error(w, "Missing or Invalid token", http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(auth, "Bearer ")
        if tokenString == "" {
            http.Error(w, "Missing or Invalid token", http.StatusUnauthorized)
            return
        }

        claims := &m.JWT{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, http.ErrNotSupported
            }
            return signingKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        if claims.Exp < time.Now().Unix() {
            http.Error(w, "Token Expired", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "username", claims.Username)
        next.ServeHTTP(w, r.WithContext(ctx))

    })
}
