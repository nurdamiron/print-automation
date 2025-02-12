// internal/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "time"
    "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Пропускаем авторизацию для некоторых эндпоинтов
            if r.URL.Path == "/api/v1/auth/login" || r.URL.Path == "/api/v1/auth/register" {
                next.ServeHTTP(w, r)
                return
            }

            // Получаем токен из заголовка
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            // Проверяем формат Bearer token
            tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
            if tokenString == authHeader {
                http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
                return
            }

            // Парсим и проверяем токен
            claims := &Claims{}
            token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(jwtSecret), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            // Добавляем информацию о пользователе в контекст
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            ctx = context.WithValue(ctx, "email", claims.Email)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}