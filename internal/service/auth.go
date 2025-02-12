// internal/service/auth.go
package service

import (
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
    "print-automation/internal/models"
    "print-automation/internal/repository"
)

type AuthService struct {
    userRepo  *repository.UserRepository
    jwtSecret string
}

// Изменяем тип UserID на string в Claims
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
    }
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
    exists, err := s.userRepo.ExistsByEmail(email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("user already exists")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user, err := s.userRepo.Create(email, string(hashedPassword))
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
    user, err := s.userRepo.GetByEmail(email)
    if err != nil {
        return "", err
    }
    if user == nil {
        return "", errors.New("user not found")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    if err != nil {
        return "", errors.New("invalid password")
    }

    claims := &Claims{
        UserID: user.ID,  // Теперь это работает, так как оба поля string
        Email:  user.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}