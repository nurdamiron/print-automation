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

            // List of public endpoints that don't require authentication
            publicEndpoints := []string{
                "/api/v1/auth/login",
                "/api/v1/auth/register",
                "/api/v1/printers/discover",
                "/api/v1/printers/status",
                "/api/v1/printers/connect",
            }

            // Check if the current path is in the public endpoints list
            for _, endpoint := range publicEndpoints {
                if strings.HasPrefix(r.URL.Path, endpoint) {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            // For all other endpoints, check for authentication
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                // Instead of returning 401, we'll allow the request but without user context
                next.ServeHTTP(w, r)
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
                // If token is invalid, continue without user context
                next.ServeHTTP(w, r)
                return
            }

            // If token is valid, add user context
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            ctx = context.WithValue(ctx, "email", claims.Email)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}