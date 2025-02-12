// internal/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "fmt"
    "github.com/golang-jwt/jwt/v4"
    "github.com/sirupsen/logrus"

)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := logrus.WithFields(logrus.Fields{
                "path":   r.URL.Path,
                "method": r.Method,
            })

            // Skip auth for public endpoints
            if r.URL.Path == "/api/v1/auth/login" || r.URL.Path == "/api/v1/auth/register" {
                logger.Debug("Skipping auth for public endpoint")
                next.ServeHTTP(w, r)
                return
            }

            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                logger.Warn("Missing authorization header")
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
            if tokenString == authHeader {
                logger.Warn("Invalid authorization format")
                http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
                return
            }

            claims := &Claims{}
            token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }
                return []byte(jwtSecret), nil
            })

            if err != nil || !token.Valid {
                logger.WithError(err).Warn("Invalid token")
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            logger = logger.WithField("user_id", claims.UserID)
            logger.Debug("User authenticated successfully")

            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            ctx = context.WithValue(ctx, "email", claims.Email)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}